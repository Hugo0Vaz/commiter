{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true; # For some IDE components if needed
        };

        # Common Go version (change as needed)
        goVersion = pkgs.go_1_24;

        # Development tools
        devTools = with pkgs; [
          # Core Go toolchain
          goVersion
          gopls            # Official Go language server
          delve            # Debugger
          go-tools         # Static analysis (staticcheck, etc)
          golangci-lint    # Popular linter aggregator
          gomodifytags     # Go struct tags manipulation
          gotests          # Generate tests
          richgo           # Colorized `go test` output
          mockgen          # Generate mocks

          # Additional utilities
          air              # Live reload for Go apps
          gofumpt          # Strict gofmt
          revive           # Fast linter
          govulncheck      # Vulnerability scanner
          sqlc            # SQL to Go type-safe codegen
          protobuf        # Protocol Buffers
          grpcurl         # gRPC debugging

          # Optional IDE (VSCode with Go extensions)
          (vscode-with-extensions.override {
            vscode = vscodium; # or vscode
            vscodeExtensions = with vscode-extensions; [
              golang.go
              ms-vscode.makefile-tools
            ];
          })
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          name = "go-dev-env";

          packages = devTools;

          # Environment variables
          GOPATH = "${toString ./.}/.go";
          GOBIN = "${toString ./.}/.go/bin";
          GO111MODULE = "on";

          shellHook = ''
            export PATH="$GOBIN:$PATH"
            echo "Go `${pkgs.go}/bin/go version`"
            echo "Available tools:"
            echo "  gopls, delve, golangci-lint, air, sqlc, etc."
          '';
        };

        # For projects using Go modules
        apps.${system} = {
          run = {
            type = "app";
            program = "${goVersion}/bin/go";
          };
        };
      });
}
