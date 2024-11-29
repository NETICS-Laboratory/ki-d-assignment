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
import { useNavigate } from "react-router-dom";

function Upload() {
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    file: null,
  });

  const handleChange = (e) => {
    const { files } = e.target;
    setFormData({
      file: files[0], // Get the first file from the file input
    });
  };

  const handleFileUpload = async (e) => {
    e.preventDefault();
    const data = new FormData(); // FormData for file upload
    Object.keys(formData).forEach((key) => {
      data.append(key, formData[key]);
    });

    try {
      const response = await axios.post("http://127.0.0.1:8090/api/files/upload", data, {
        headers: {
          "Content-Type": "multipart/form-data", // Set correct headers for file upload
        },
      });

      console.log(response);

      // Redirect user after successful registration
      navigate("/upload");
    } catch (error) {
      console.error("Error signing up:", error.response?.data || error.message);
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
                Upload File
              </SoftTypography>
              <SoftTypography variant="button" fontWeight="regular">
                Fill in the form to upload a file
              </SoftTypography>
            </SoftBox>
            <Separator />
            <SoftBox pt={2} pb={3} px={3}>
              <SoftBox component="form" role="form" onSubmit={handleFileUpload}>
                <SoftBox mb={2}>
                  <SoftInput
                    type="file"
                    name="id_card"
                    onChange={handleChange}
                    // accept="image/*" // Optional: limit to image files
                  />
                </SoftBox>

                <SoftBox mt={4} mb={1}>
                  <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                    Upload
                  </SoftButton>
                </SoftBox>
              </SoftBox>
            </SoftBox>
          </Card>
        </SoftBox>
      </SoftBox>
    </DashboardLayout>
  );
}

export default Upload;
