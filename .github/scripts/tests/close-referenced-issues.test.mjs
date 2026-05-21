// Copyright 2026 The MathWorks, Inc.
import { describe, it, expect, vi } from "vitest";
import { createRequire } from "module";

const require = createRequire(import.meta.url);
const closeReferencedIssues = require("../src/close-referenced-issues.js");
const { extractIssueNumbers, closeIssue } = closeReferencedIssues;

describe("extractIssueNumbers", () => {
    it("extracts issue numbers from an Issues Resolved section", () => {
        const body = [
            "## What's New",
            "Some features",
            "## Issues Resolved",
            "- Fixed #42",
            "- Fixed #99",
            "## Other",
            "- #200 should not be included",
        ].join("\n");

        expect(extractIssueNumbers(body)).toEqual([42, 99]);
    });

    it("returns empty array when no Issues Resolved section exists", () => {
        const body = "## What's New\n- Added feature #10";
        expect(extractIssueNumbers(body)).toEqual([]);
    });

    it("deduplicates issue numbers", () => {
        const body = "## Issues Resolved\n- #5 and #5 again";
        expect(extractIssueNumbers(body)).toEqual([5]);
    });

    it("handles empty body", () => {
        expect(extractIssueNumbers("")).toEqual([]);
    });

    it("is case insensitive for section heading", () => {
        const body = "## issues resolved\n- #7";
        expect(extractIssueNumbers(body)).toEqual([7]);
    });

    it("stops at next heading of same or higher level", () => {
        const body = [
            "## Issues Resolved",
            "- #1",
            "## Next Section",
            "- #2",
        ].join("\n");

        expect(extractIssueNumbers(body)).toEqual([1]);
    });

    it("ignores cross-repo references like org/repo#123", () => {
        const body = [
            "## Issues Resolved",
            "- Fixed #10",
            "- See also mathworks/other-repo#234",
            "- Related to org/repo#456",
        ].join("\n");

        expect(extractIssueNumbers(body)).toEqual([10]);
    });
});

describe("closeIssue", () => {
    function createMocks({ issueState = "open", isPR = false } = {}) {
        return {
            github: {
                rest: {
                    issues: {
                        get: vi.fn().mockResolvedValue({
                            data: {
                                state: issueState,
                                pull_request: isPR ? {} : undefined,
                            },
                        }),
                        createComment: vi.fn().mockResolvedValue({}),
                        update: vi.fn().mockResolvedValue({}),
                    },
                },
            },
            core: {
                info: vi.fn(),
                warning: vi.fn(),
            },
        };
    }

    it("closes an open issue and adds a comment", async () => {
        const { github, core } = createMocks();

        await closeIssue({
            github,
            owner: "org",
            repo: "repo",
            issueNumber: 42,
            tag: "v1.0.0",
            releaseUrl: "https://github.com/org/repo/releases/tag/v1.0.0",
            core,
        });

        expect(github.rest.issues.createComment).toHaveBeenCalledWith(
            expect.objectContaining({ issue_number: 42 }),
        );
        expect(github.rest.issues.update).toHaveBeenCalledWith(
            expect.objectContaining({ state: "closed", state_reason: "completed" }),
        );
    });

    it("skips already closed issues", async () => {
        const { github, core } = createMocks({ issueState: "closed" });

        await closeIssue({
            github,
            owner: "org",
            repo: "repo",
            issueNumber: 42,
            tag: "v1.0.0",
            releaseUrl: "https://example.com",
            core,
        });

        expect(github.rest.issues.createComment).not.toHaveBeenCalled();
        expect(github.rest.issues.update).not.toHaveBeenCalled();
        expect(core.info).toHaveBeenCalledWith(expect.stringContaining("already closed"));
    });

    it("skips pull requests", async () => {
        const { github, core } = createMocks({ isPR: true });

        await closeIssue({
            github,
            owner: "org",
            repo: "repo",
            issueNumber: 42,
            tag: "v1.0.0",
            releaseUrl: "https://example.com",
            core,
        });

        expect(github.rest.issues.createComment).not.toHaveBeenCalled();
        expect(github.rest.issues.update).not.toHaveBeenCalled();
        expect(core.info).toHaveBeenCalledWith(expect.stringContaining("pull request"));
    });
});

describe("closeReferencedIssues", () => {
    it("logs info and returns when no issues found", async () => {
        const core = { info: vi.fn(), warning: vi.fn() };
        const context = {
            payload: { release: { body: "No issues here", tag_name: "v1.0.0", html_url: "" } },
            repo: { owner: "org", repo: "repo" },
        };

        await closeReferencedIssues({ github: {}, context, core });

        expect(core.info).toHaveBeenCalledWith("No issue references found in release notes.");
    });

    it("handles null release body", async () => {
        const core = { info: vi.fn(), warning: vi.fn() };
        const context = {
            payload: { release: { body: null, tag_name: "v1.0.0", html_url: "" } },
            repo: { owner: "org", repo: "repo" },
        };

        await closeReferencedIssues({ github: {}, context, core });

        expect(core.info).toHaveBeenCalledWith("No issue references found in release notes.");
    });

    it("handles non-404 errors gracefully", async () => {
        const error = new Error("Internal Server Error");
        error.status = 500;

        const github = {
            rest: { issues: { get: vi.fn().mockRejectedValue(error) } },
        };
        const core = { info: vi.fn(), warning: vi.fn() };
        const context = {
            payload: {
                release: {
                    body: "## Issues Resolved\n- #123",
                    tag_name: "v1.0.0",
                    html_url: "",
                },
            },
            repo: { owner: "org", repo: "repo" },
        };

        await closeReferencedIssues({ github, context, core });

        expect(core.warning).toHaveBeenCalledWith(
            expect.stringContaining("Failed to close #123"),
        );
    });

    it("handles 404 errors gracefully", async () => {
        const error = new Error("Not Found");
        error.status = 404;

        const github = {
            rest: { issues: { get: vi.fn().mockRejectedValue(error) } },
        };
        const core = { info: vi.fn(), warning: vi.fn() };
        const context = {
            payload: {
                release: {
                    body: "## Issues Resolved\n- #999",
                    tag_name: "v1.0.0",
                    html_url: "",
                },
            },
            repo: { owner: "org", repo: "repo" },
        };

        await closeReferencedIssues({ github, context, core });

        expect(core.warning).toHaveBeenCalledWith(expect.stringContaining("#999 not found"));
    });
});
