import { useEffect, useState } from "react";
import { QRCode } from "react-qr-svg";
import axios from "axios";
import { makeClaimRequest} from "../utils/utils";
import { useRouter } from "next/router";
import { Layout } from "../components";
import { Flex, Heading, Paragraph } from "theme-ui";

const Page = (props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
  const [loading, setLoading] = useState(true);
  const [qrData, setQRData] = useState({});

  const checkVerificationStatus = async (sessionID: string) => {
    try {
      const resp = await axios.get(
          "http://" +props.issuerLocalUrl + `/api/v1/status?id=${sessionID}`
      );
      if (resp.status === 200) {
        return true;
      }

      return false;

    } catch (err) {
      // TODO: Error Handling
      console.log("err->", err);
      return false;
    }
  };

  useEffect(() => {
    (async () => {
      const resp = await axios.get(
          "http://" +props.issuerLocalUrl + "/api/v1/age-verification-request"
      );

      setQRData(resp.data);
      setLoading(false);

      const sessionID = resp.headers["x-id"];

      const interval = setInterval(async () => {
        const isVerified = await checkVerificationStatus(sessionID);
        if (isVerified) {
          clearInterval(interval);
          alert("verification succeeded ✅");
        }
      }, 2000);
    })();
  }, []);
  return (
    <Layout>
      {loading ? (
        <h1>Loading</h1>
      ) : (
        <Flex
          sx={{ flex: 1, flexDirection: "column", variant: "layout.allCenter" }}
        >
          <Heading sx={{ textAlign: "center", fontSize: [32], my: [4] }}>
            Verify your claim 👮‍♀️
          </Heading>

          <QRCode
            level="Q"
            style={{ width: 256 }}
            value={JSON.stringify(qrData)}
          />

          <Paragraph sx={{ variant: "text.para" }}>
            Scan this to verify you are above 22 years old.
          </Paragraph>
        </Flex>
      )}
    </Layout>
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
