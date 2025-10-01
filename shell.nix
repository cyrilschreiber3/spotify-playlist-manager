{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
          templ.overlays.default
        ];
      }
  ),
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
  templ ? pkgs.templ,
}: let
  goEnv = mkGoEnv {pwd = ./.;};

  prettier-plugin-go-template-patched = pkgs.prettier-plugin-go-template.overrideAttrs (oldAttrs: {
    postInstall =
      oldAttrs.postInstall
      + ''
        if [[ $nodeModulesPath == *prettier-plugin-go-template/node_modules ]]; then
          echo "Fixing node modules location"
          mv $nodeModulesPath/* $out/lib/node_modules/
        fi
      '';
  });

  prettierConfig = pkgs.writeTextFile {
    name = "prettierrc";
    text = ''
      {
        "plugins": [
          "${prettier-plugin-go-template-patched}/lib/node_modules/prettier-plugin-go-template/lib/index.js"
        ],
        "overrides": [
          {
            "files": "*.gohtml",
            "options": {
              "parser": "go-template"
            }
          }
        ]
      }
    '';
  };
in
  with pkgs;
    mkShell {
      packages = [
        # gomod2nix prerequisites
        goEnv
        gomod2nix

        # Go development
        air
        delve
        go
        golangci-lint
        golangci-lint-langserver
        gomodifytags
        gopls
        gotests
        impl

        # Templ
        templ

        # Database
        sqlc

        # Web development
        nodePackages.prettier
        prettier-plugin-go-template-patched
        tailwindcss_4
      ];

      shellHook = ''
        ln -sf ${prettierConfig} ./.prettierrc
        sed -i '/.*prettier.prettierPath*/c\  "prettier.prettierPath\": "${prettier-plugin-go-template-patched}/lib/node_modules/prettier",' ./.vscode/settings.json

        echo -e "Welcome to the Go dev environment!\n"

        echo -e "$(${go}/bin/go version)\n"

      '';
    }
