import type { AppProps } from "next/app";
import { ThemeProvider, Image } from "theme-ui";
import theme from "../utils/theme";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <ThemeProvider theme={theme}>
      <Component {...pageProps} />
      <Image
        src="./logo.png"
        sx={{ position: "fixed", width: [240], right: [3], bottom: [3] }}
      />
    </ThemeProvider>
  );
}

export default MyApp;
