# Binary Package Manager (bpm)

Binary Package Manager keeps track of binary tools installed from separate (GitHub) repos. Besides just knowing about them, `bpm` can also install, update and remove them.

## Why?

Some ecosystems (hello K8s) consist of multiple little tools that are downloadable as binaries from their specific repos. Some of them notify you when there's a new version available, but keeping them up to date is quite cumbersome. `bpm` tries to solve this by giving you yet another binary so instead of all others you'll need to keep track of just one. Going forward, the plan is to make `bpm` available via Homebrew, Debian packaging etc so you can keep it up to date via other means.
