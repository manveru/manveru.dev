{
  description = "Flake for the Ganso-Nisshodo website";

  inputs = {
    nixpkgs.follows = "euphenix/nixpkgs";
    flake-utils.follows = "euphenix/flake-utils";
    euphenix.url = "gitlab:manveru/euphenix/flake";
  };

  outputs = { self, euphenix, flake-utils, nixpkgs }:
    (flake-utils.lib.simpleFlake {
      inherit self nixpkgs;
      name = "site";
      preOverlays = [ euphenix.overlay ];
      overlay = ./overlay.nix;
      shell = ./shell.nix;
    }) // {
      overlay = import ./overlay.nix;
    };
}
