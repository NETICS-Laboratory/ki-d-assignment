import React, { useState, useEffect } from "react";
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

import axios from "axios";

function DecryptedFile() {
  const [id, setId] = useState("");
  const [isDecrypted, setIsDecrypted] = useState(false);
  const [rows, setRows] = useState([]);
  const [downloadPath, setDownloadPath] = useState("");

  // Table columns
  const columns = [
    { name: "data", label: "Data", align: "center" },
    { name: "plaintext", label: "Plaintext", align: "left" },
  ];

  // Handle decrypt button click
  const handleDecrypt = async (e) => {
    e.preventDefault();

    try {
      const token = localStorage.getItem("token");

      const response = await axios.post(
        "http://127.0.0.1:8090/api/files/get-file-decrypted",
        { id }, // Send file ID in the request body
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const decryptedData = response.data.data;

      setRows([
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Decrypted AES
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {decryptedData.decrypted_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Decrypted RC4
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {decryptedData.decrypted_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Decrypted DES
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {decryptedData.decrypted_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Signature
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {decryptedData.signature}
            </SoftTypography>
          ),
        },
      ]);

      // Modify the path by changing "encrypted" to "decrypted" and removing ".aes"
      const modifiedPath = decryptedData.decrypted_aes
        .replace("encrypted", "decrypted") // Replace "encrypted" with "decrypted"
        .replace(".aes", ""); // Remove the ".aes" part

      // Store the modified path for download
      setDownloadPath(modifiedPath);

      // Store file paths for download buttons
      // setDownloadPath(decryptedData.decrypted_aes);

      setIsDecrypted(true);
    } catch (error) {
      console.error("Error during decryption:", error);
      alert("Decryption failed. Please try again.");
    }
  };

  // Download handler
  const handleDownload = (path) => {
    // Trigger file download by opening the file URL
    const link = document.createElement("a");
    link.href = `http://127.0.0.1:8090/${path}`; // Append the path to your base API URL
    link.download = path.split("/").pop(); // Use the file name from the path
    link.click();
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <SoftBox mb={3}>
          {!isDecrypted ? (
            <Card>
              <SoftBox p={3} textAlign="center">
                <SoftTypography variant="h5" fontWeight="medium">
                  Decrypt File
                </SoftTypography>
                <SoftTypography variant="button" fontWeight="regular">
                  Enter the ID of the file you want to decrypt
                </SoftTypography>
              </SoftBox>
              <Separator />
              <SoftBox pt={2} pb={3} px={3}>
                <SoftBox component="form" role="form" onSubmit={handleDecrypt}>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="text"
                      name="id"
                      placeholder="File ID"
                      value={id}
                      onChange={(e) => setId(e.target.value)}
                    />
                  </SoftBox>

                  <SoftBox mt={4} mb={1}>
                    <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                      Decrypt
                    </SoftButton>
                  </SoftBox>
                </SoftBox>
              </SoftBox>
            </Card>
          ) : (
            <SoftBox>
              <Card>
                <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
                  <SoftTypography variant="h6">Decrypted Data</SoftTypography>
                </SoftBox>
                <SoftBox
                  sx={{
                    "& .MuiTableRow-root:not(:last-child)": {
                      "& td": {
                        borderBottom: ({ borders: { borderWidth, borderColor } }) =>
                          `${borderWidth[1]} solid ${borderColor}`,
                      },
                    },
                  }}
                >
                  <Table columns={columns} rows={rows} />
                </SoftBox>
                <SoftBox mt={3} p={3}>
                  <SoftButton
                    variant="gradient"
                    color="info"
                    fullWidth
                    onClick={() => handleDownload(downloadPath)}
                  >
                    Download File
                  </SoftButton>
                </SoftBox>
              </Card>
            </SoftBox>
          )}
        </SoftBox>
      </SoftBox>
    </DashboardLayout>
  );
}

export default DecryptedFile;
