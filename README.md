# Kompresi

[日本語訳はこちら](https://github.com/tsg0o0/Kompresi/blob/master/README_JA.md)

Go application that compresses PNG and JPEG images in folders in the background with no degradation.

**This application is still in alpha testing.
In particular, I do not have a lot of development knowledge about Windows.
Any assistance is welcome!**

## What can this do?

Monitors the specified directory and compresses any detected image files losslessly.

[Zopfli](https://github.com/google/zopfli) is used for PNG compression and [Guetzli](https://github.com/google/guetzli) for JPEG compression.

Zopfli and Guetzli are **very slow** to compress.
This application is intended to run in the background and is not suitable for anything else.

## Download & Setup

### 1. Download

Download the latest version [here](https://github.com/tsg0o0/Kompresi/releases).

*You can also build from source code if you have Go installed on your PC.*

### 2. Complete the configuration

**For Windows and macOS, the included KompresiConfigure (KompresiConfigure.exe) allows for easy on-screen configuration and startup. **

Open a terminal (command prompt) and run kompresi (kompresi.exe).
You will probably see a setup guide.

The following arguments are used to configure the settings.

- `inputDir 'YOUR INPUT DIRECTRY PATH'` Select the directory to load the images.
- `outputDir 'YOUR OUTPUT DIRECTRY PATH'` Select a directory to output compressed images.
- `deleteOrigin 'Yes or No'` Delete original files after compression.
- `optimLv '0 - 2'` Change the compression level.
  - `0`: Fast but low compression
  - `1`: Auto (experimental)
  - `2`: Slow but high compression
- `help` Show help.
- `license` Show lisence.

### 3. Run

Run again with no arguments. Starts running in the background.

<sub> *If a specific image file is specified as an argument, compression of that image is initiated. This is for performance testing.* </sub>

## Support

If you find a bug or have a question, please create an [Issue](https://github.com/tsg0o0/Kompresi/issues) or [contact me](https://tsg0o0.com/contact/).

## License

This software is licensed under the terms of the [Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/).

## Tip

Did you like it? [Send me a tip](https://tsg0o0.com/tip/) if you like!
