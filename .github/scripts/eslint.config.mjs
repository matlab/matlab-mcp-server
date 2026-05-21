// Copyright 2026 The MathWorks, Inc.
import js from "@eslint/js";
import nodePlugin from "eslint-plugin-n";

export default [
    {
        ignores: ["eslint.config.mjs", "tests/"],
    },
    js.configs.recommended,
    nodePlugin.configs["flat/recommended-script"],
    {
        languageOptions: {
            ecmaVersion: "latest",
            sourceType: "commonjs",
            globals: {
                require: "readonly",
                module: "readonly",
                console: "readonly",
                process: "readonly",
            },
        },
        rules: {
            curly: ["error", "all"],
            eqeqeq: ["error", "always"],
            "no-var": "error",
            "prefer-const": "error",
            "no-throw-literal": "error",
            "no-return-await": "error",
            "no-shadow": "error",
            "no-param-reassign": "error",
            "no-implicit-coercion": "error",
            "no-else-return": "error",
            "no-template-curly-in-string": "warn",
            "no-useless-concat": "error",
            "prefer-template": "error",
            "object-shorthand": ["error", "always"],
            "dot-notation": "warn",
            "array-callback-return": "warn",
            "no-unreachable-loop": "warn",
            "no-void": "warn",
            "prefer-regex-literals": "warn",
        },
    },
];
