{
  description = "flemmy: stuff and flemmys";

  inputs.nixpkgs.url = "nixpkgs/nixos-23.05";

  outputs = { self, nixpkgs }:
    let
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in {
      overlay = _: prev: { inherit (self.packages.${prev.system}) flemmy; };

      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          flemmy = pkgs.buildGoModule {
            pname = "flemmy";
            version = "v0.0.0";
            src = ./.;

            vendorHash = "sha256-id4nO3w8qmwmn4FJBgs/E19MtbjMjJ0fDqiFYsHIIsY=";
            proxyVendor = true;

            nativeBuildInputs = with pkgs; [ pkg-config ];
            buildInputs = with pkgs; [
              git
              glfw
              libGL
              libGLU
              openssh
              pkg-config
              xorg.libXcursor
              xorg.libXi
              xorg.libXinerama
              xorg.libXrandr
              xorg.libXxf86vm
              xorg.xinput
            ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.flemmy);
      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            shellHook = ''
              PS1='\u@\h:\@; '
              nix flake run github:qbit/xin#flake-warn
              echo "Go `${pkgs.go}/bin/go version`"
            '';
            buildInputs = with pkgs; [
              git
              go_1_20
              gopls
              go-tools

              glfw
              pkg-config
              xorg.libXcursor
              xorg.libXi
              xorg.libXinerama
              xorg.libXrandr
              xorg.libXxf86vm
              xorg.xinput
            ];
          };
        });
    };
}

