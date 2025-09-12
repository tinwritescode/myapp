import { Box, Button, Heading, Text, VStack } from "@chakra-ui/react";
import { useState } from "react";
import { LoginForm } from "./components/LoginForm";
import { RegisterForm } from "./components/RegisterForm";
import { useAuthStore } from "./hooks/useAuthStore";

function App() {
  const { isLoggedIn, user, logout } = useAuthStore();
  const [isLoginMode, setIsLoginMode] = useState(true);

  const handleSwitchToRegister = () => setIsLoginMode(false);
  const handleSwitchToLogin = () => setIsLoginMode(true);

  return (
    <Box minH="100vh" bg="gray.50" pt={8}>
      {isLoggedIn ? (
        <Box p={8}>
          <VStack gap={4} align="center">
            <Heading size="lg">Welcome back!</Heading>
            <Text>You are logged in as: {user?.email}</Text>
            <Text>Full name: {user?.full_name}</Text>
            <Text>Username: {user?.username}</Text>
            <Button colorPalette="red" onClick={logout}>
              Logout
            </Button>
          </VStack>
        </Box>
      ) : isLoginMode ? (
        <LoginForm onSwitchToRegister={handleSwitchToRegister} />
      ) : (
        <RegisterForm onSwitchToLogin={handleSwitchToLogin} />
      )}
    </Box>
  );
}

export default App;
