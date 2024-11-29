import React, { useState } from "react";
import Card from "@mui/material/Card";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

// Soft UI Dashboard React components
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftInput from "components/SoftInput";
import SoftButton from "components/SoftButton";

// Authentication layout components
import Separator from "layouts/authentication/components/Separator";

// Soft UI Dashboard React examples
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import Table from "examples/Tables/Table";

// Axios for API calls
import axios from "axios";

// Material UI Snackbar for notifications
import { Snackbar, Alert } from "@mui/material";

function VerifySignature() {
  const [file_id, setFileID] = useState("");
  const [signature, setSignature] = useState("");
  const [openSnackbar, setOpenSnackbar] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [snackbarSeverity, setSnackbarSeverity] = useState("success"); // "success" or "error"

  const handleVerify = async (e) => {
    e.preventDefault();

    const token = localStorage.getItem("token");

    try {
      // Call the verify signature API
      const response = await axios.post(
        "http://127.0.0.1:8090/api/files/verify-digital-signature",
        {
          file_id,
          signature,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      // If verification is successful
      setSnackbarMessage("Verification Successful!");
      setSnackbarSeverity("success");
      setOpenSnackbar(true);
    } catch (error) {
      console.error("Error during verifying signature:", error);

      // If there was an error in verification
      setSnackbarMessage("Verification Failed! Please try again.");
      setSnackbarSeverity("error");
      setOpenSnackbar(true);
    }
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <SoftBox mb={3}>
          <Card>
            <SoftBox p={3} textAlign="center">
              <SoftTypography variant="h5" fontWeight="medium">
                Verify Digital Signature
              </SoftTypography>
              <SoftTypography variant="button" fontWeight="regular">
                Fill in the form to verify the digital signature
              </SoftTypography>
            </SoftBox>
            <Separator />
            <SoftBox pt={2} pb={3} px={3}>
              <SoftBox component="form" role="form" onSubmit={handleVerify}>
                <SoftBox mb={2}>
                  <SoftInput
                    type="text"
                    name="file_id"
                    placeholder="File ID"
                    value={file_id}
                    onChange={(e) => setFileID(e.target.value)}
                  />
                </SoftBox>
                <SoftBox mb={2}>
                  <SoftInput
                    type="text"
                    name="signature"
                    placeholder="Signature"
                    value={signature}
                    onChange={(e) => setSignature(e.target.value)}
                  />
                </SoftBox>

                <SoftBox mt={4} mb={1}>
                  <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                    Verify
                  </SoftButton>
                </SoftBox>
              </SoftBox>
            </SoftBox>
          </Card>
        </SoftBox>
      </SoftBox>

      {/* Snackbar for showing success/error message */}
      <Snackbar
        open={openSnackbar}
        autoHideDuration={6000}
        onClose={() => setOpenSnackbar(false)}
        anchorOrigin={{
          vertical: "bottom", // Position the Snackbar at the top
          horizontal: "center", // Center it horizontally
        }}
      >
        <Alert
          onClose={() => setOpenSnackbar(false)}
          severity={snackbarSeverity}
          sx={{ width: "100%" }}
        >
          {snackbarMessage}
        </Alert>
      </Snackbar>
    </DashboardLayout>
  );
}

export default VerifySignature;
