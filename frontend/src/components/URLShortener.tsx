import {
  Box,
  Button,
  Heading,
  HStack,
  Input,
  Text,
  VStack,
  Spinner,
  Badge,
  IconButton,
} from "@chakra-ui/react";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { urlApi } from "@/lib/urlApi";
import type { CreateURLRequest } from "@/types/url";
import { toaster } from "./ui/toaster";

export function URLShortener() {
  const [form, setForm] = useState<CreateURLRequest>({
    original_url: "",
    short_code: "",
    expires_at: "",
  });
  const [createdURL, setCreatedURL] = useState<string | null>(null);

  // Create URL mutation
  const createURLMutation = useMutation({
    mutationFn: (data: CreateURLRequest) => urlApi.createPublicURL(data),
    onSuccess: (response) => {
      setCreatedURL(response.data.short_code);
      setForm({ original_url: "", short_code: "", expires_at: "" });
      toaster.create({
        title: "URL shortened successfully!",
        type: "success",
        duration: 3000,
      });
    },
    onError: (error: unknown) => {
      const errorResponse = error as {
        response?: { data?: { error?: string; code?: string } };
      };
      const errorCode = errorResponse?.response?.data?.code;
      const errorMessage = errorResponse?.response?.data?.error;

      if (errorCode === "SHORT_CODE_ALREADY_EXISTS") {
        toaster.create({
          title: "Short code already exists",
          description:
            "Please choose a different short code or leave it empty for auto-generation",
          type: "error",
          duration: 5000,
        });
      } else {
        toaster.create({
          title: "Failed to shorten URL",
          description:
            errorMessage ||
            (error as Error)?.message ||
            "An unexpected error occurred",
          type: "error",
          duration: 5000,
        });
      }
    },
  });

  const handleSubmit = () => {
    if (!form.original_url.trim()) {
      toaster.create({
        title: "URL is required",
        type: "error",
        duration: 3000,
      });
      return;
    }

    const data: CreateURLRequest = {
      original_url: form.original_url,
      short_code: form.short_code || undefined,
      expires_at: form.expires_at || undefined,
    };

    createURLMutation.mutate(data);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toaster.create({
      title: "Copied to clipboard",
      type: "success",
      duration: 2000,
    });
  };

  const getShortURL = (shortCode: string) => {
    const baseURL =
      import.meta.env.VITE_API_URL?.replace("/api/v1", "") ||
      "https://myapp-1757744589.fly.dev";
    return `${baseURL}/${shortCode}`;
  };

  return (
    <Box p={6}>
      <VStack gap={6} align="stretch" maxW="600px" mx="auto">
        <Box p={6} bg="white" borderRadius="md" shadow="sm">
          <VStack gap={4}>
            <Heading size="lg" textAlign="center">
              URL Shortener
            </Heading>
            <Text textAlign="center" color="gray.600">
              Create short, memorable links for your URLs
            </Text>

            <VStack gap={4} w="full">
              <Box w="full">
                <Text mb={2}>Original URL *</Text>
                <Input
                  placeholder="https://example.com/very/long/url"
                  value={form.original_url}
                  onChange={(e) =>
                    setForm({ ...form, original_url: e.target.value })
                  }
                  size="lg"
                />
              </Box>

              <Box w="full">
                <Text mb={2}>Custom Short Code (optional)</Text>
                <Input
                  placeholder="mycode"
                  value={form.short_code}
                  onChange={(e) =>
                    setForm({ ...form, short_code: e.target.value })
                  }
                />
                <Text fontSize="sm" color="gray.600" mt={1}>
                  6-8 characters, alphanumeric only
                </Text>
              </Box>

              <Box w="full">
                <Text mb={2}>Expiration Date (optional)</Text>
                <Input
                  type="datetime-local"
                  value={form.expires_at}
                  onChange={(e) =>
                    setForm({ ...form, expires_at: e.target.value })
                  }
                />
              </Box>

              <Button
                colorPalette="blue"
                size="lg"
                onClick={handleSubmit}
                loading={createURLMutation.isPending}
                w="full"
              >
                {createURLMutation.isPending ? (
                  <HStack>
                    <Spinner size="sm" />
                    <Text>Shortening...</Text>
                  </HStack>
                ) : (
                  "Shorten URL"
                )}
              </Button>
            </VStack>
          </VStack>
        </Box>

        {/* Result */}
        {createdURL && (
          <Box p={6} bg="white" borderRadius="md" shadow="sm">
            <VStack gap={4}>
              <Heading size="md" textAlign="center">
                Your Short URL
              </Heading>
              <Box
                w="full"
                p={4}
                bg="green.50"
                borderRadius="md"
                border="1px solid"
                borderColor="green.200"
              >
                <HStack justify="space-between">
                  <VStack align="start" gap={1}>
                    <Text fontFamily="mono" fontSize="lg" fontWeight="bold">
                      {getShortURL(createdURL)}
                    </Text>
                    <Badge colorPalette="green">Ready to use</Badge>
                  </VStack>
                  <IconButton
                    variant="outline"
                    onClick={() => copyToClipboard(getShortURL(createdURL))}
                  >
                    📋
                  </IconButton>
                </HStack>
              </Box>

              <Button
                variant="outline"
                onClick={() => setCreatedURL(null)}
                size="sm"
              >
                Create Another
              </Button>
            </VStack>
          </Box>
        )}
      </VStack>
    </Box>
  );
}
