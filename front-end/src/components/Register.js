import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { API_URL } from "../Env";
import "../styles/style.css";
import "../styles/Register.css";

const Register = () => {
  const [username, setUsername] = useState("");
  const [email_address, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordConfirm, setPasswordConfirm] = useState("");
  const [birthday, setBirthday] = useState("01");
  const [birthmonth, setBirthmonth] = useState("01");
  const [birthyear, setBirthyear] = useState(
    new Date().getFullYear().toString()
  );
  const [birthdate, setBirthdate] = useState("01-01-2023");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    setBirthdate(`${birthyear}-${birthmonth}-${birthday}`);
  }, [birthday, birthmonth, birthyear]);

  const handleKeyDown = (event) => {
    if (event.key === "Enter") {
      handleRegister();
    }
  };

  const handleRegister = async () => {
    if (password !== passwordConfirm) {
      setError("Password and password confirmation do not match");
      return;
    }

    try {
      const response = await fetch(API_URL + "/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email_address, password, birthdate }),
      });
      if (response.status === 201) {
        const data = await response.json();
        console.log(data);
      } else {
        const returnMessage = await response.json();
        const errorMessage =
          returnMessage.error.charAt(0).toUpperCase() +
          returnMessage.error.slice(1);
        console.log(
          JSON.stringify({ username, email_address, password, birthdate })
        );
        console.error(returnMessage);
        setError(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError("An error occurred while registering, please try again later");
    }
  };

  return (
    <div
      className="background-color"
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        height: "100vh",
      }}
    >
      <h1 className="text" style={{ marginTop: "3rem" }}>
        Create an account
      </h1>
      {error && (
        <div className="error_box" style={{ marginBottom: "1rem" }}>
          {error}
        </div>
      )}
      <div
        style={{
          width: "30%",
          display: "flex",
          flexDirection: "column",
          alignItems: "flex-start",
        }}
      >
        <label
          className="text"
          style={{ alignSelf: "flex-start", marginBottom: "0.5rem" }}
        >
          Username
        </label>
        <input
          type="text"
          id="username"
          className="text_box"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <label
          className="text"
          style={{
            alignSelf: "flex-start",
            marginTop: "0.5rem",
            marginBottom: "0.5rem",
          }}
        >
          Email
        </label>
        <input
          type="email"
          id="email"
          className="text_box"
          value={email_address}
          onChange={(e) => setEmail(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <label
          className="text"
          style={{
            alignSelf: "flex-start",
            marginTop: "0.5rem",
            marginBottom: "0.5rem",
          }}
        >
          Password
        </label>
        <input
          type="password"
          id="password"
          className="text_box"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <label
          className="text"
          style={{
            alignSelf: "flex-start",
            marginTop: "0.5rem",
            marginBottom: "0.5rem",
          }}
        >
          Confirm Password
        </label>
        <input
          type="password"
          id="passwordConfirm"
          className="text_box"
          value={passwordConfirm}
          onChange={(e) => setPasswordConfirm(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <label
          className="text"
          style={{
            alignSelf: "flex-start",
            marginTop: "0.5rem",
            marginBottom: "0.5rem",
          }}
        >
          Birth date
        </label>
        <div style={{ display: "flex", alignItems: "center" }}>
          <select
            id="birthday"
            className="text_box"
            value={birthday}
            onChange={(e) => setBirthday(e.target.value)}
            onKeyDown={handleKeyDown}
            style={{ marginRight: "0.5rem", width: "7.5rem" }}
          >
            {Array.from(Array(31), (_, i) =>
              (i + 1).toString().padStart(2, "0")
            ).map((day) => (
              <option key={day} value={day}>
                {day}
              </option>
            ))}
          </select>
          <select
            id="birthmonth"
            className="text_box"
            value={birthmonth}
            onChange={(e) => setBirthmonth(e.target.value)}
            onKeyDown={handleKeyDown}
            style={{ marginRight: "0.5rem", width: "7.5rem" }}
          >
            {Array.from(Array(12), (_, i) =>
              (i + 1).toString().padStart(2, "0")
            ).map((month) => (
              <option key={month} value={month}>
                {month}
              </option>
            ))}
          </select>
          <select
            id="birthyear"
            className="text_box"
            value={birthyear}
            onChange={(e) => setBirthyear(e.target.value)}
            onKeyDown={handleKeyDown}
            style={{ width: "7.5rem" }}
          >
            {Array.from(Array(100), (_, i) => new Date().getFullYear() - i).map(
              (year) => (
                <option key={year} value={year}>
                  {year}
                </option>
              )
            )}
          </select>
        </div>
        <button
          className="button"
          style={{
            padding: "1rem",
            marginTop: "2rem",
            alignSelf: "center",
            fontFamily: "Fira Code",
            fontSize: "1rem",
          }}
          onClick={handleRegister}
        >
          Register
        </button>
      </div>
    </div>
  );
};

export default Register;
