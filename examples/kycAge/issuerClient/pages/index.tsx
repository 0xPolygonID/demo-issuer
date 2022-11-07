import { useEffect, useState } from "react";
import axios from "axios";
import { checkAuthStatus } from "../utils/utils";
import { useRouter } from "next/router";
import { Container, Flex, Heading, Paragraph, Spinner } from "theme-ui";
import { Layout, QRCode } from "../components";

const Page = (props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
  const [loading, setLoading] = useState(true);
  const [qrData, setQRData] = useState({});

  console.log("frotend log, Issuer Public Url:", props.issuerPublicUrl);
  console.log("frotend log, Issuer Local Url:", props.issuerLocalUrl);

  const router = useRouter();

  useEffect(() => {
    (async () => {
      const resp = await axios.get(props.issuerLocalUrl + "/api/sign-in");

      setQRData(resp.data);
      setLoading(false);

      const sessionID = resp.headers["x-id"];

      const interval = setInterval(async () => {
        const resp = await checkAuthStatus(sessionID, props);
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


export async function getServerSideProps(context) {
  const yaml = require('js-yaml');
  const fs = require('fs');

  let  issuerPublicUrl = "";
  let  issuerLocalUrl = "";


  try {
    const doc = yaml.load(fs.readFileSync('./../../../issuer/issuer_config.default.yaml', 'utf8'));
    issuerPublicUrl = doc.public_url;
    issuerLocalUrl = doc.local_url;
  } catch (e) {
    console.log("encounter error on load config file, err: " + e);
    process.exit(1);
  }

  return {
    props: {issuerPublicUrl: issuerPublicUrl, issuerLocalUrl: issuerLocalUrl },
  }
}

export default Page;
