let
  # euphenixSource = import ~/github/manveru/euphenix { };
  euphenixSource = (import (fetchTarball {
    url =
      "https://github.com/manveru/euphenix/archive/34252d3ec764e91d8578728573fa8350ce5127fe.tar.gz";
  }) { });

  euphenix = euphenixSource.extend (self: super: {
    parseMarkdown =
      super.parseMarkdown.override { flags = { prismjs = true; }; };
  });

  inherit (euphenix.lib) take optionalString nameValuePair hasPrefix;
  inherit (euphenix) build sortByRecent loadPosts pp mkPostCSS cssTag;

  withLayout = body: [ ./templates/layout.html body ];
  variables = rec {
    activeClass = route: prefix:
      optionalString (hasPrefix prefix route) "active";
    posts = sortByRecent (loadPosts "/blog/" ./blog);
    latestPosts = take 10 posts;
    liveJS = ''<script src="/js/live.js"></script>'';
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
  };

  blogRoutes = __listToAttrs (map (post:
    nameValuePair post.url {
      template = withLayout ./templates/blog_show.html;
      variables = variables // post // { title = post.meta.title; };
    }) variables.posts);

in euphenix.build ./. { routes = pageRoutes // blogRoutes; }
