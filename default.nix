{ euphenix }:
let
  inherit (euphenix.lib)
    take optionalString nameValuePair hasPrefix elemAt length;
  inherit (euphenix) build sortByRecent loadPosts pp mkPostCSS cssTag;

  withLayout = body: [ ./templates/layout.html body ];
  variables = rec {
    activeClass = route: prefix:
      optionalString (hasPrefix prefix route) "active";
    posts = sortByRecent (loadPosts "/blog/" ./blog);
    latestPosts = take 10 posts;
    liveJS = optionalString ((__getEnv "LIVEJS") != "")
      ''<script src="/js/live.js"></script>'';
    css = cssTag (mkPostCSS ./css);
  };

  route = template: title: {
    template = withLayout template;
    variables = variables // { inherit title; };
  };

  pageRoutes = {
    "/index.html" = route ./templates/home.html "Home";
    "/contact/index.html" = route ./templates/contact.html "Contact";
    "/404/index.html" = route ./templates/404.html "Not Found";
    "/archive/index.html" = route ./templates/archive.html "Blog Archive";
    "/blog/index.html" = route ./templates/blog.html "Blog";
    "/blog/feed.atom" = {
      template = [ ./templates/atom.xml ];
      variables = variables // {
        updated =
          (elemAt variables.posts (length variables.posts - 1)).meta.date;
        author = "Michael 'manveru' Fellinger";
      };
    };
  };

  blogRoutes = __listToAttrs (map (post:
    nameValuePair post.url {
      template = withLayout ./templates/blog_show.html;
      variables = variables // post // { title = post.meta.title; };
    }) variables.posts);

in euphenix.build {
  src = ./.;
  routes = pageRoutes // blogRoutes;
}
