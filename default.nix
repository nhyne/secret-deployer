{ 
  pkgs ? import (fetchTarball https://github.com/NixOS/nixpkgs/archive/79b15cdf49e2174f50f0384b1f8188538469ac03.tar.gz) {}
}:
with pkgs;
let
    go = pkgs.go_1_15;
    buildInputs = [
      go
    ];
in {
    shell = pkgs.mkShell {
        buildInputs = buildInputs;
    };
}
