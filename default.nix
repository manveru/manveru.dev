let
  pkgs = import ./nix/nixpkgs.nix;
  euphenix = import ~/github/manveru/euphenix { };
  inherit (pkgs.lib) take;
  inherit (euphenix) build loadPosts sortByRecent;
in build {
  rootDir = ./.;
  layout = ./templates/layout.tmpl;
  favicon = ./static/img/favicon.svg;

  variables = rec {
    posts = sortByRecent (loadPosts "/blog/" ./content/posts);
    latestPosts = take 5 posts;
  };
}
