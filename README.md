<p align="center">
  <h1 align="center"> IdaThemer 🎨</h1>
  <p align="center">
    <a href="https://github.com/can1357/IdaThemer/actions/workflows/go.yml">
      <img alt="github-actions" src="https://img.shields.io/github/actions/workflow/status/can1357/IdaThemer/go.yml"/>
    </a>
    <a href="https://github.com/can1357/IdaThemer/blob/master/LICENSE.md">
      <img alt="license" src="https://img.shields.io/github/license/can1357/IdaThemer.svg"/>
    </a>
  </p>

  <p align="center">
      Seamlessly convert your favorite Visual Studio Code themes to IDA Pro themes.
  </p>
</p>

![Result](https://i.can.ac/YjBan.png)

## Introduction

IdaThemer is a tool for converting Visual Studio Code (VsCode) themes into compatible themes for IDA Pro. This utility leverages the [Dracula](https://github.com/dracula/ida) theme as its base and then remaps the colors defined in the source and target VsCode theme JSON files to create a visually appealing and functional theme for IDA Pro.

This project is still in its early stages and may produce unexpected results. If you encounter any issues, please open an issue on GitHub with the output of the program and the VsCode theme you are trying to convert. Although light themes _should be_ supported, the primary focus of this project is to provide a dark theme for IDA Pro. As such, light themes may not be as polished as dark themes.

You can find some of the most popular themes in the [themes](https://github.com/can1357/IdaThemer/tree/master/themes) directory.

## Building from source

Make sure you have [Go > 1.21](https://golang.org/) installed on your system. Then run the following commands:

```bash
git clone https://github.com/can1357/IdaThemer.git
cd IdaThemer
go build
```

## Converting a single theme

Run IdaThemer with the path to your VsCode theme JSON and the output directory:

```bash
./IdaThemer "gruvbox.json" "S:\IDA Pro\themes"
```

## Converting all themes from VSCode

Run IdaThemer with \* as the first argument and the output directory as the second argument:

```bash
./IdaThemer * "S:\IDA Pro\themes"
```

If VSCode is not installed at `$HOME` or `%USERPROFILE%`, you should set the Environment Variable `VSCODE_DATA` to the path of your VSCode configuration where `.vscode/extensions` is located before running IdaThemer.

When the process is complete, you should have all of your themes converted and ready to use in IDA Pro.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgments

Thanks to the Zeno Rocha and the Dracula team for their work on their IDA Pro theme: [Dracula](https://github.com/dracula/ida)
