// Copyright 2026 The MathWorks, Inc.

package installmatlab

const (
	name  = "install_matlab"
	title = "Install MATLAB"
	description = `Downloads and installs MATLAB and/or MATLAB toolboxes using the MATLAB Package Manager (mpm).
This tool does NOT require MATLAB to already be installed. It downloads the mpm CLI and uses it to install the specified products.
Product names use underscores for spaces (e.g., Signal_Processing_Toolbox, Deep_Learning_Toolbox).
To install MATLAB itself, include "MATLAB" in the products list.`
)

// Args defines the input schema for the install_matlab tool.
type Args struct {
	Release     string   `json:"release"     jsonschema:"The MATLAB release to install (e.g. R2025a or R2024b)."`
	Products    []string `json:"products"    jsonschema:"List of product names to install. Use underscores for spaces (e.g. MATLAB or Signal_Processing_Toolbox or Deep_Learning_Toolbox)."`
	Destination string   `json:"destination,omitempty" jsonschema:"The installation directory. If not specified mpm uses its default location."`
}

// ReturnArgs defines the output schema for the install_matlab tool.
type ReturnArgs struct {
	Output string `json:"output" jsonschema:"The output from the mpm install command."`
}
