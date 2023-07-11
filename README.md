# Sprites 🧚

A CLI tool that generates SVG spritesheets from a directory of SVG files.

> ⚠️ Please don't use this package, this is a personal experimental project.

## Installation

```sh
npm install @brattonross/sprites
```

## Usage

Given the following directory structure:

```
public
└── svg
    └── icons
        ├── 20
        │   └── solid
        │       ├── arrow-left-on-rectangle.svg
        │       ├── pencil-square.svg
        │       └── user.svg
        └── 24
            └── outline
                ├── computer-desktop.svg
                ├── moon.svg
                └── sun.svg
```

Running the following command:

```sh
sprites --src public/svg/icons --sprites public/svg/spritesheets --components src/components/icons
```

Will generate the following structure:

```
public
└── svg
    ├── icons
    │   ├── 20
    │   │   └── solid
    │   │       ├── arrow-left-on-rectangle.svg
    │   │       ├── pencil-square.svg
    │   │       └── user.svg
    │   └── 24
    │       └── outline
    │           ├── computer-desktop.svg
    │           ├── moon.svg
    │           └── sun.svg
    └── spritesheets
        ├── 20-solid.svg
        └── 24-outline.svg
src
└── components
    └── icons
        ├── 20-solid.tsx
        └── 24-outline.tsx
```

It is assumed that you will be serving the generated spritesheets from the `public` directory.

## Prior Art

- [The "best" way to manage icons in React.js](https://benadam.me/thoughts/react-svg-sprites/)
- [rmx-cli](https://github.com/kiliman/rmx-cli#-svg-sprite--new)
