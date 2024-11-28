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

function RequestData() {
  const navigate = useNavigate();
  const [requested_username, setRequested_username] = useState("");

  const handleSendRequest = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem("token");

      const response = await axios.post(
        "http://127.0.0.1:8090/api/user/request-access",
        {
          requested_username,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`, // Add the Bearer token here
          },
        }
      );

      // Reload the page after a successful request
      window.location.reload();
    } catch (error) {
      console.error("Error sending request:", error.response?.data || error.message);
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
                Request Data
              </SoftTypography>
              <SoftTypography variant="button" fontWeight="regular">
                Fill in the form to request access to other user&apos;s data
              </SoftTypography>
            </SoftBox>
            <Separator />
            <SoftBox pt={2} pb={3} px={3}>
              <SoftBox component="form" role="form" onSubmit={handleSendRequest}>
                <SoftBox mb={2}>
                  <SoftInput
                    type="text"
                    placeholder="Username"
                    value={requested_username}
                    onChange={(e) => setRequested_username(e.target.value)}
                  />
                </SoftBox>

                <SoftBox mt={4} mb={1}>
                  <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                    Send Request
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

export default RequestData;
