import { useEffect, useState } from "react";
import axios from "axios";
import { checkAuthStatus } from "../utils/utils";
import { useRouter } from "next/router";
import { Container, Flex, Heading, Paragraph, Spinner } from "theme-ui";
import { Layout, QRCode } from "../components";

const Page = () => {
  const [loading, setLoading] = useState(true);
  const [qrData, setQRData] = useState({});

  const router = useRouter();

  useEffect(() => {
    (async () => {
      const resp = await axios.get("http://localhost:3000/api/sign-in");

      setQRData(resp.data);
      setLoading(false);

      const sessionID = resp.headers["x-id"];

      const interval = setInterval(async () => {
        const resp = await checkAuthStatus(sessionID);
        if (resp) {
          clearInterval(interval);
          router.push(`/client?claimID=${resp.claimID}&userID=${resp.userID}`);
        }
      }, 2000);
    })();
  }, []);

  return (
    <Container>
      {loading ? (
        <>
          <Heading>Loading</Heading>
          <Spinner color="purple" />
        </>
      ) : (
        <Layout>
          <Flex
            sx={{
              flex: 1,
              justifyContent: "center",
              alignItems: "center",
              flexDirection: "column",
            }}
          >
            <QRCode
              level="Q"
              style={{ width: "100%" }}
              value={JSON.stringify(qrData)}
            />
            <Paragraph sx={{ variant: "text.para" }}>
              Scan the QR code to sign in with Polygon ID
            </Paragraph>
          </Flex>
        </Layout>
      )}
    </Container>
  );
};

export default Page;
