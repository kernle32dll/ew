# ew

ew - short for `(run things) e(very)w(here)` is a tool for grouping folders by tags,
and executing tasks in all folders via these tags.

## How to use

COMMAND                           | DESCRIPTION
-------                           | ----
ew                                |   list all paths, grouped by their tags
ew help                           |   displays this help
ew --help                         |   displays this help (alias for ew help)
ew migrate                        |   migrate from mixu/gr config, and keep json format
ew migrate --yaml                 |   migrate from mixu/gr config, and use new yaml format
ew paths                          |   list all paths (alias for ew paths list)
ew paths list                     |   list all paths
ew tags                           |   list all tags (alias for ew tags list)
ew tags list                      |   list all tags
ew tags add @some-tag             |   add current directory to tag "some-tag" (NOT YET IMPLEMENTED)
ew tags add \some\path @some-tag  |   add \some\path to tag "some-tag" (NOT YET IMPLEMENTED)
ew tags rm @some-tag              |   add current directory to tag "some-tag" (NOT YET IMPLEMENTED)
ew tags rm \some\path @some-tag   |   add \some\path to tag "some-tag" (NOT YET IMPLEMENTED)
ew status                         |   show quick git status for all paths
ew @tag1 status                   |   show quick git status for all paths of tag1 (supports multiple tags)
ew @tag1 some-cmd                 |   executes some-cmd in all paths of tag1 (supports multiple tags)
