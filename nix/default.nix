{ sources ? import ./sources.nix }:
with {
  overlay = _: pkgs: {
    inherit (import sources.niv { }) niv;
    packages = pkgs.callPackages ./packages.nix { };
    euphenix = (import sources.euphenix { }).extend (self: super: {
      parseMarkdown =
        super.parseMarkdown.override { flags = { prismjs = true; }; };
    });
  };
};
import sources.nixpkgs {
  overlays = [ overlay ];
  config = { };
}
