let
  euphenixSource = import ~/github/manveru/euphenix { };
  # euphenixSource = (import (fetchTarball {
  #   url =
  #     "https://github.com/manveru/euphenix/archive/eaaee37df12e8fccced3a4ac93402d6b3e5fcf54.tar.gz";
  # }) { })

  euphenix = euphenixSource.extend (self: super: {
    parseMarkdown =
      super.parseMarkdown.override { flags = { prismjs = true; }; };
  });
  inherit (euphenix.lib) take optionalString;
  inherit (euphenix) build sortByRecent loadPosts pp;

  mkRoutes = args:
    builtins.attrValues (builtins.mapAttrs (route: value:
      let file = builtins.toFile (builtins.baseNameOf route) value.body;
      in euphenix.mkDerivation {
        name = "mkRoute-${builtins.baseNameOf route}";
        buildInputs = [ euphenix.coreutils ];
        inherit route;
        buildCommand = ''
          install -m 0644 -D ${file} "$out${route}"
        '';
      }) args);

in euphenix.build {
  rootDir = ./.;
  layout = ./templates/layout.html;
  favicon = ./static/img/favicon.svg;

  routes = mkRoutes {
    "/index.html".body = scopedImport {
      title = "Home";
      cssTag = "";
      activeClass = expected: optionalString ("/" == expected) "active";
      route = "/";
      liveJS = ''<script src="/js/live.js"></script>'';
      content = scopedImport {
        latestPosts = take 10 (sortByRecent (loadPosts "/blog/" ./blog));
      } ./templates/home.html;
    } ./templates/layout.html;

    "/contact/index.html" = {
      # inputs = [(css ./css/main.css)];
      body = scopedImport {
        title = "Contact";
        cssTag = "";
        activeClass = expected:
          optionalString ("/contact/index.html" == expected) "active";
        route = "/contact/index.html";
        liveJS = ''<script src="/js/live.js"></script>'';
        content = import ./templates/contact.html;
      } ./templates/layout.html;
    };
  };

  variables = rec {
    liveJS = ''<script src="/js/live.js"></script>'';
    inherit pp;
    posts = sortByRecent (loadPosts "/blog/" ./blog);
    latestPosts = take 20 posts;
  };

  expensiveVariables = rec {
    posts = sortByRecent (loadPosts "/blog/" ./blog);
    latestPosts = take 20 posts;
  };
}
