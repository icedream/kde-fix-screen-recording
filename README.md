
# KDE screen recording fix tool

This little tool is a hack to get screen capturing working again on KDE.

## Why?

The KDE compositor is known to cause various issues with screen recording, such as tearing, screen flickering[^1] and in my case eye-blindingly fast flashing between frames captured in the past and the present. The only known way to work around this is to disable screen effects by disabling said compositor.

## Build

```bash
go build -v github.com/icedream/kde-fix-screen-recording
```

## License

This project is licensed under the MIT license. For more information check the [`LICENSE`](LICENSE) file.

