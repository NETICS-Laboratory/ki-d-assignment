import React, { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import Table from "examples/Tables/Table";
import axios from "axios";

function EncryptedFile() {
  const [data, setData] = useState([]);

  // Fetch data function
  const fetchData = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.get("http://127.0.0.1:8090/api/files/get-files", {
        headers: {
          Authorization: `Bearer ${token}`, // Add the Bearer token here
        },
      });

      const responseData = response.data.data;
      setData(responseData);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    fetchData(); // Call the fetch function on component load
  }, []);

  const renderTable = (file) => {
    const columns = [
      { name: "data", align: "center", width: "20%" },
      { name: "ciphertext", align: "left", width: "80%" },
    ];

    const rows = [
      {
        data: (
          <SoftTypography variant="caption" fontWeight="medium">
            ID
          </SoftTypography>
        ),
        ciphertext: (
          <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            {file.id}
          </SoftTypography>
        ),
      },
      {
        data: (
          <SoftTypography variant="caption" fontWeight="medium">
            File AES
          </SoftTypography>
        ),
        ciphertext: (
          <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            {file.files_aes}
          </SoftTypography>
        ),
      },
      {
        data: (
          <SoftTypography variant="caption" fontWeight="medium">
            File RC4
          </SoftTypography>
        ),
        ciphertext: (
          <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            {file.files_rc4}
          </SoftTypography>
        ),
      },
      {
        data: (
          <SoftTypography variant="caption" fontWeight="medium">
            File DES
          </SoftTypography>
        ),
        ciphertext: (
          <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            {file.files_des}
          </SoftTypography>
        ),
      },
      {
        data: (
          <SoftTypography variant="caption" fontWeight="medium">
            Signature
          </SoftTypography>
        ),
        ciphertext: (
          <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            {file.signature}
          </SoftTypography>
        ),
      },
    ];

    return (
      <SoftBox mb={3} key={file.id}>
        <Card>
          <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
            <SoftTypography variant="h6">File</SoftTypography>
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
        </Card>
      </SoftBox>
    );
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        {data.length > 0 ? (
          data.map((file) => renderTable(file))
        ) : (
          <SoftTypography variant="caption" fontWeight="medium">
            No files found.
          </SoftTypography>
        )}
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default EncryptedFile;
