let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  unstable = import sources.nixpkgs_unstable { };
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
    unstable.golangci-lint # use unstable to get version of golangci-lint that works with go 1.14
    pkgs.gotools
    pkgs.mage
    pkgs.nodejs-12_x
    pkgs.yarn
  ];
}
