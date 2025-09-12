import {
  Alert,
  Box,
  Button,
  Field,
  Input,
  Text,
  VStack,
} from "@chakra-ui/react";
import { useState } from "react";
import { useAuthStore } from "@/hooks/useAuthStore";

interface LoginFormProps {
  onSuccess?: () => void;
  onSwitchToRegister?: () => void;
}

export const LoginForm = ({
  onSuccess,
  onSwitchToRegister,
}: LoginFormProps) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuthStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      await login(email, password);
      onSuccess?.();
    } catch (err: unknown) {
      let errorMessage = "Login failed. Please try again.";

      if (err instanceof Error && "response" in err) {
        const response = (
          err as { response?: { data?: { error?: string; code?: string } } }
        ).response;
        const errorData = response?.data;

        if (errorData?.code === "INVALID_CREDENTIALS") {
          errorMessage =
            "Invalid email or password. Please check your credentials and try again.";
        } else if (errorData?.code === "UNAUTHORIZED") {
          errorMessage =
            "Your account has been deactivated. Please contact support.";
        } else if (errorData?.error) {
          errorMessage = errorData.error;
        }
      }

      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box
      maxW="md"
      mx="auto"
      pt={8}
      px={6}
      pb={6}
      borderWidth={1}
      borderRadius="lg"
    >
      <Text fontSize="xl" fontWeight="bold" mb={6} textAlign="center">
        Sign In
      </Text>

      <form onSubmit={handleSubmit}>
        <VStack gap={4}>
          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Email
            </Text>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
            />
          </Field.Root>

          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Password
            </Text>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
            />
          </Field.Root>

          {error && (
            <Alert.Root status="error">
              <Alert.Indicator />
              <Alert.Content>
                <Alert.Title>Error!</Alert.Title>
                <Alert.Description>{error}</Alert.Description>
              </Alert.Content>
            </Alert.Root>
          )}

          <Button
            type="submit"
            colorPalette="blue"
            width="full"
            loading={isLoading}
            loadingText="Signing in..."
          >
            Sign In
          </Button>

          <Text fontSize="sm" textAlign="center" mt={2}>
            Don't have an account?{" "}
            <Button
              variant="ghost"
              colorPalette="blue"
              onClick={onSwitchToRegister}
              p={0}
              h="auto"
              fontSize="sm"
            >
              Sign up
            </Button>
          </Text>
        </VStack>
      </form>
    </Box>
  );
};
