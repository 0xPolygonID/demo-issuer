import { FunctionComponent } from "react";
import { QRCode as QrCode, QRCodeProps } from "react-qr-svg";
import { Box } from "theme-ui";

const QRCode: FunctionComponent<QRCodeProps & React.SVGProps<SVGElement>> = (
  props
) => {
  return (
    <Box sx={{ width: [320] }}>
      <QrCode level="Q" style={{ width: "100%" }} {...props} />
    </Box>
  );
};

export default QRCode;
