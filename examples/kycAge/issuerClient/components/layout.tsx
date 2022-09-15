import { FunctionComponent, PropsWithChildren } from "react";
import { Flex } from "theme-ui";

const Layout: FunctionComponent<PropsWithChildren> = ({ children }) => {
  return (
    <Flex
      sx={{
        height: "100vh",
        width: "100vw",
        flexDirection: "column",
        backgroundImage: `url( "./bg.png")`,
        backgroundSize: "cover",
        backgroundPosition: "center",
      }}
    >
      {children}
    </Flex>
  );
};

export default Layout;
