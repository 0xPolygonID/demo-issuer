import { useEffect, useState } from "react";
import axios from "axios";
import {makeAgeClaimRequest} from "../utils/utils";
import { useRouter } from "next/router";
import { Container, Flex, Heading, Paragraph, Spinner } from "theme-ui";
import { Layout, QRCode } from "../components";

const Page = (props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
  const [loading, setLoading] = useState(true);
  const [qrData, setQRData] = useState({});

  const router = useRouter();

  useEffect(() => {
    (async () => {

      const resp = await axios.get("http://" + props.issuerLocalUrl + "/api/v1/requests/auth");

      setQRData(resp.data);
      setLoading(false);

      const sessionID = resp.headers["x-id"];

      const interval = setInterval(async () => {
        try {

          const resp = await axios.get("http://" + props.issuerLocalUrl + `/api/v1/status?id=${sessionID}`);
          console.log('Here', resp.data);

          if (resp) {
            const userID = resp.data.id
            if (userID) {
              let dob = 19860503;
              const respMakeClaim = await axios(makeAgeClaimRequest(dob, userID, props));
              // TODO: Error Handling
              const claimID = respMakeClaim.data.id ? respMakeClaim.data.id : "";
              clearInterval(interval);
              router.push(`/client?claimID=${claimID}&userID=${userID}&dob=${dob}`);
            }
          }
        } catch (e) {
          console.log('err->', e);
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

  let issuerPublicUrl = "";
  let issuerLocalUrl = "";

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
