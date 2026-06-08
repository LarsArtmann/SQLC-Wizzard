{
  description = "SQLC-Wizard — Interactive CLI wizard for generating sqlc configurations";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    systems.url = "github:nix-systems/default";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          config,
          pkgs,
          lib,
          ...
        }:
        let
          version = self.rev or self.dirtyRev or "dev";
          commit = self.rev or "unknown";
          buildDate = self.lastModifiedDate or "1970-01-01T00:00:00Z";

          vendorHash = "sha256-8YQp5yy3Xb/7TOd2klsOAeml8VNA4rCxjJSRAa/VYVY=";
          proxyVendor = true;

          ldflags = [
            "-s"
            "-w"
            "-X main.Version=${version}"
            "-X main.Commit=${commit}"
            "-X main.BuildDate=${buildDate}"
          ];

          sqlc-wizard =
            {
              withTests ? false,
            }:
            pkgs.buildGoModule {
              pname = "sqlc-wizard";
              inherit version vendorHash;

              src = lib.fileset.toSource {
                root = ./.;
                fileset = lib.fileset.unions [
                  ./go.mod
                  ./go.sum
                  ./cmd
                  ./internal
                  ./pkg
                  ./generated
                ];
              };

              subPackages = [ "cmd/sqlc-wizard" ];

              inherit ldflags;

              env.CGO_ENABLED = 0;

              tags = [
                "netgo"
                "osusergo"
              ];

              doCheck = withTests;

              checkFlags =
                if withTests then
                  [
                    "-v"
                    "-race"
                  ]
                else
                  [ ];

              meta = with lib; {
                description = "Interactive CLI wizard for generating sqlc configurations";
                homepage = "https://github.com/LarsArtmann/SQLC-Wizzard";
                license = licenses.mit;
                mainProgram = "sqlc-wizard";
                maintainers = [ maintainers.larsartmann ];
                platforms = platforms.unix;
              };
            };
        in
        {
          packages = {
            default = sqlc-wizard { };
          };

          devShells = {
            default = pkgs.mkShell {
              packages =
                with pkgs;
                [
                  go
                  gopls
                  golangci-lint
                  goreleaser
                  bun
                  typescript
                  go-arch-lint
                ]
                ++ lib.optionals stdenv.isLinux [ gdb ];

              inputsFrom = [ config.packages.default ];

              GOWORK = "off";

              shellHook = ''
                echo "sqlc-wizard dev shell"
                echo "  go: $(go version)"
                echo "  golangci-lint: $(golangci-lint version --format short 2>/dev/null || echo 'unknown')"
              '';
            };

            ci = pkgs.mkShellNoCC {
              packages = with pkgs; [
                go
                golangci-lint
              ];
            };          };

          apps.default = {
            type = "app";
            program = lib.getExe config.packages.default;
          };

          checks = {
            format = config.treefmt.build.check self;
            build = config.packages.default;

            test = sqlc-wizard { withTests = true; };

            vet =
              pkgs.runCommandLocal "sqlc-wizard-vet"
                {
                  src = lib.fileset.toSource {
                    root = ./.;
                    fileset = lib.fileset.unions [
                      ./go.mod
                      ./go.sum
                      ./cmd
                      ./internal
                      ./pkg
                      ./generated
                    ];
                  };
                  nativeBuildInputs = [ pkgs.go_1_26 ];
                }
                ''
                  cd $src
                  go vet ./...
                  touch $out
                '';

            gofmt =
              pkgs.runCommandLocal "sqlc-wizard-gofmt"
                {
                  src = lib.fileset.toSource {
                    root = ./.;
                    fileset = lib.fileset.unions [
                      ./cmd
                      ./internal
                      ./pkg
                      ./generated
                    ];
                  };
                  nativeBuildInputs = [ pkgs.go_1_26 ];
                }
                ''
                  cd $src
                  test -z "$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)
                  touch $out
                '';
          };

          treefmt.config = {
            projectRootFile = "go.mod";
            programs = {
              nixfmt.enable = true;
              gofmt.enable = true;
              gofumpt.enable = true;
              goimports.enable = true;
            };
          };
        };

      flake.overlays.default = final: _prev: {
        SQLC-Wizzard = self.packages.${final.stdenv.system}.default;
      };

    };
}
