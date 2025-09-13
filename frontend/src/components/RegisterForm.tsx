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
import { Link } from "react-router-dom";
import { useAuthStore } from "@/hooks/useAuthStore";

interface RegisterFormProps {
  onSuccess?: () => void;
}

export const RegisterForm = ({ onSuccess }: RegisterFormProps) => {
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
    confirmPassword: "",
    fullName: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { register } = useAuthStore();

  const handleInputChange =
    (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({ ...prev, [field]: e.target.value }));
    };

  const validateForm = () => {
    if (formData.password !== formData.confirmPassword) {
      setError("Passwords do not match");
      return false;
    }
    if (formData.password.length < 6) {
      setError("Password must be at least 6 characters long");
      return false;
    }
    if (formData.username.length < 3) {
      setError("Username must be at least 3 characters long");
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      await register({
        email: formData.email,
        username: formData.username,
        password: formData.password,
        full_name: formData.fullName,
      });
      onSuccess?.();
    } catch (err: unknown) {
      let errorMessage = "Registration failed. Please try again.";

      if (err instanceof Error && "response" in err) {
        const response = (
          err as { response?: { data?: { error?: string; code?: string } } }
        ).response;
        const errorData = response?.data;

        if (errorData?.code === "EMAIL_ALREADY_USED") {
          errorMessage =
            "This email or username is already registered. Please use a different email or username.";
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
    <Box maxW="md" mx="auto" mt={8} p={6} borderWidth={1} borderRadius="lg">
      <Text fontSize="xl" fontWeight="bold" mb={6} textAlign="center">
        Create Account
      </Text>

      <form onSubmit={handleSubmit}>
        <VStack gap={4}>
          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Full Name
            </Text>
            <Input
              type="text"
              value={formData.fullName}
              onChange={handleInputChange("fullName")}
              placeholder="Enter your full name"
            />
          </Field.Root>

          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Email
            </Text>
            <Input
              type="email"
              value={formData.email}
              onChange={handleInputChange("email")}
              placeholder="Enter your email"
            />
          </Field.Root>

          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Username
            </Text>
            <Input
              type="text"
              value={formData.username}
              onChange={handleInputChange("username")}
              placeholder="Choose a username"
            />
          </Field.Root>

          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Password
            </Text>
            <Input
              type="password"
              value={formData.password}
              onChange={handleInputChange("password")}
              placeholder="Enter your password"
            />
          </Field.Root>

          <Field.Root required>
            <Text fontSize="sm" fontWeight="medium" mb={1}>
              Confirm Password
            </Text>
            <Input
              type="password"
              value={formData.confirmPassword}
              onChange={handleInputChange("confirmPassword")}
              placeholder="Confirm your password"
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
            loadingText="Creating account..."
          >
            Create Account
          </Button>

          <Text fontSize="sm" textAlign="center" mt={2}>
            Already have an account?{" "}
            <Link to="/login">
              <Button
                variant="ghost"
                colorPalette="blue"
                p={0}
                h="auto"
                fontSize="sm"
              >
                Sign in
              </Button>
            </Link>
          </Text>
        </VStack>
      </form>
    </Box>
  );
};
