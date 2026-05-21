// Copyright 2026 The MathWorks, Inc.

package custom

import "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"

func (t *Tool) SetToolAdder(toolAdder basetool.ToolAdder[map[string]any, any]) {
	t.toolAdder = toolAdder
}
