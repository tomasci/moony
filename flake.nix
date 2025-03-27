{
  description = "Moony Server Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        myInputs = [
          # go deps
          pkgs.go
          pkgs.gopls
          pkgs.gotools
          pkgs.golangci-lint
          pkgs.air
          pkgs.clang
          pkgs.llvm
          # sys utils
          pkgs.bash
          pkgs.coreutils
          pkgs.findutils # find command
          pkgs.which # which
          pkgs.gnused # sed
          pkgs.gnugrep # grep
          pkgs.nano # nano
          pkgs.ncurses # TUI, tput
          pkgs.docker
          pkgs.docker-compose
        ];

        # build a clean PATH from the inputs
        cleanPath = builtins.concatStringsSep ":" (map (pkg: "${pkg}/bin") myInputs);

        moonyShell = pkgs.writeTextFile {
          name = "shell.sh";
          destination = "/shell.sh";
          text = builtins.readFile ./shell.sh;
          executable = true;
        };

      in
      {
        devShells.default = pkgs.mkShell {
          pure = true;
          buildInputs = myInputs;

          shellHook = ''
            # Download Go module dependencies automatically.
            # This ensures your go.mod and go.sum are used to pull the correct versions.
            echo "Updating Go dependencies..."
            go mod download

            tput clear

            mkdir -p .nix-go

            export PS1="\[\033[01;32m\][Moony]\[\033[00m\] \w $ "

            # use clean PATH that includes only declared Nix inputs
            export PATH=${cleanPath}

            export GOPATH="$PWD/.nix-go"
            export GOBIN="$GOPATH/bin"
            export PATH="$GOBIN:$PATH"
            IDE_GO_ROOT="$(dirname $(dirname $(which go)))/share/go"

            echo ""
            echo "🌙 Moony development environment loaded"
            echo ""
            echo "Go version: $(go version)"
            echo ""
            echo "IDE configuration:"
            echo "  GOROOT: $IDE_GO_ROOT"
            echo "  GOPATH: $GOPATH"
            echo ""

            source ${moonyShell}/shell.sh
          '';
        };
      }
    );
}