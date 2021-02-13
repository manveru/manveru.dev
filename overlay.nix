final: prev: {
  site.site = prev.callPackage ./. { };
  euphenix = prev.euphenix.extend (efinal: eprev: {
    parseMarkdown =
      eprev.parseMarkdown.override { flags = { prismjs = true; }; };
  });
}
