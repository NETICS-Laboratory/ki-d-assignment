/**
=========================================================
* Soft UI Dashboard React - v4.0.1
=========================================================

* Product Page: https://www.creative-tim.com/product/soft-ui-dashboard-react
* Copyright 2023 Creative Tim (https://www.creative-tim.com)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

import axios from "axios";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

// react-router-dom components
import { Link } from "react-router-dom";

// @mui material components
import Card from "@mui/material/Card";
import Checkbox from "@mui/material/Checkbox";

// Soft UI Dashboard React components
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftInput from "components/SoftInput";
import SoftButton from "components/SoftButton";

// Authentication layout components
import BasicLayout from "layouts/authentication/components/BasicLayout";
import Socials from "layouts/authentication/components/Socials";
import Separator from "layouts/authentication/components/Separator";

// Images
import curved6 from "assets/images/curved-images/curved14.jpg";

function SignUp() {
  const navigate = useNavigate();
  const [agreement, setAgremment] = useState(true);

  const handleSetAgremment = () => setAgremment(!agreement);

  const [formData, setFormData] = useState({
    name: "",
    email: "",
    username: "",
    password: "",
    no_telp: "",
    address: "",
    id_card: null,
  });

  const handleChange = (e) => {
    const { name, value, files } = e.target;
    setFormData({
      ...formData,
      [name]: files ? files[0] : value, // For file input, get the first file from 'files'
    });
  };

  const handleSignUp = async (e) => {
    e.preventDefault();
    const data = new FormData(); // FormData for file upload
    Object.keys(formData).forEach((key) => {
      data.append(key, formData[key]);
    });

    try {
      const response = await axios.post("http://127.0.0.1:8090/api/user/register", data, {
        headers: {
          "Content-Type": "multipart/form-data", // Set correct headers for file upload
        },
      });

      console.log(response);

      // Redirect user after successful registration
      navigate("/authentication/sign-in");
    } catch (error) {
      console.error("Error signing up:", error.response?.data || error.message);
    }
  };

  return (
    <BasicLayout title="Sign up" image={curved6}>
      <Card>
        <SoftBox p={3} textAlign="center">
          <SoftTypography variant="h5" fontWeight="medium">
            Register
          </SoftTypography>
          <SoftTypography variant="button" fontWeight="regular">
            Fill in the data below to register
          </SoftTypography>
        </SoftBox>
        <Separator />
        <SoftBox pt={2} pb={3} px={3}>
          <SoftBox component="form" role="form" onSubmit={handleSignUp}>
            <SoftBox mb={2}>
              <SoftInput
                type="text"
                name="name"
                placeholder="Name"
                value={formData.name}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="email"
                name="email"
                placeholder="Email"
                value={formData.email}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="text"
                name="username"
                placeholder="Username"
                value={formData.username}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="password"
                name="password"
                placeholder="Password"
                value={formData.password}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="text"
                name="no_telp"
                placeholder="Phone Number"
                value={formData.no_telp}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="text"
                name="address"
                placeholder="Address"
                value={formData.address}
                onChange={handleChange}
              />
            </SoftBox>
            <SoftBox mb={2}>
              <SoftInput
                type="file"
                name="id_card"
                onChange={handleChange}
                accept="image/*" // Optional: limit to image files
              />
            </SoftBox>
            <SoftBox mt={4} mb={1}>
              <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                sign up
              </SoftButton>
            </SoftBox>
            <SoftBox mt={3} textAlign="center">
              <SoftTypography variant="button" color="text" fontWeight="regular">
                Already have an account?&nbsp;
                <SoftTypography
                  component={Link}
                  to="/authentication/sign-in"
                  variant="button"
                  color="dark"
                  fontWeight="bold"
                  textGradient
                >
                  Sign in
                </SoftTypography>
              </SoftTypography>
            </SoftBox>
          </SoftBox>
        </SoftBox>
      </Card>
    </BasicLayout>
  );
}

export default SignUp;
