let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
in pkgs.mkShell {
  buildInputs = [
    # Dependencies for running grafana-server in a rootless container.
    pkgs.podman
    pkgs.conmon
    pkgs.runc
    pkgs.slirp4netns
    pkgs.fuse-overlayfs

    # Developer tooling.
    pkgs.go
    pkgs.golangci-lint
    pkgs.gotools
    pkgs.mage
    pkgs.nodejs-12_x
    pkgs.yarn
  ];
}
