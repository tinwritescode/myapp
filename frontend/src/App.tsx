import { Box } from "@chakra-ui/react";
import { Navigate, Route, Routes } from "react-router-dom";
import { LoginForm } from "./components/LoginForm";
import { RegisterForm } from "./components/RegisterForm";
import { URLManager } from "./components/URLManager";
import { URLShortener } from "./components/URLShortener";
import { useAuthStore } from "./hooks/useAuthStore";

function App() {
  const { isLoggedIn } = useAuthStore();

  return (
    <Box minH="100vh">
      {isLoggedIn ? (
        <Box>
          {/* Main Content */}
          <Routes>
            <Route path="/" element={<Navigate to="/shorten" replace />} />
            <Route path="/shorten" element={<URLShortener />} />
            <Route path="/manage" element={<URLManager />} />
          </Routes>
        </Box>
      ) : (
        <Box pt={8}>
          <Routes>
            <Route path="/" element={<LoginForm />} />
            <Route path="/login" element={<LoginForm />} />
            <Route path="/register" element={<RegisterForm />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        </Box>
      )}
    </Box>
  );
}

export default App;
