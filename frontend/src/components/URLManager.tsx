import {
  Box,
  Button,
  Heading,
  HStack,
  Input,
  Text,
  VStack,
  Badge,
  IconButton,
  Spinner,
  Center,
} from "@chakra-ui/react";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { urlApi } from "@/lib/urlApi";
import type { URL, CreateURLRequest, UpdateURLRequest } from "@/types/url";
import { useAuthStore } from "@/hooks/useAuthStore";
import { toaster } from "./ui/toaster";

export function URLManager() {
  const { isLoggedIn } = useAuthStore();
  const queryClient = useQueryClient();

  // State for forms
  const [selectedURL, setSelectedURL] = useState<URL | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showEditForm, setShowEditForm] = useState(false);

  // Form states
  const [createForm, setCreateForm] = useState<CreateURLRequest>({
    original_url: "",
    short_code: "",
    expires_at: "",
  });
  const [editForm, setEditForm] = useState<UpdateURLRequest>({
    original_url: "",
    expires_at: "",
    is_active: true,
  });

  // Fetch URLs
  const {
    data: urlsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["urls", page, limit, searchTerm],
    queryFn: () =>
      urlApi.getURLs({
        page,
        limit,
        search: searchTerm || undefined,
        sort_by: "created_at",
        sort_dir: "desc",
      }),
    enabled: isLoggedIn,
  });

  // Create URL mutation
  const createURLMutation = useMutation({
    mutationFn: (data: CreateURLRequest) =>
      isLoggedIn ? urlApi.createURL(data) : urlApi.createPublicURL(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["urls"] });
      setShowCreateForm(false);
      setCreateForm({ original_url: "", short_code: "", expires_at: "" });
      toaster.create({
        title: "URL created successfully",
        type: "success",
        duration: 3000,
      });
    },
    onError: (error: unknown) => {
      toaster.create({
        title: "Failed to create URL",
        description:
          (error as { response?: { data?: { message?: string } } })?.response
            ?.data?.message || (error as Error)?.message,
        type: "error",
        duration: 5000,
      });
    },
  });

  // Update URL mutation
  const updateURLMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateURLRequest }) =>
      urlApi.updateURL(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["urls"] });
      setShowEditForm(false);
      setSelectedURL(null);
      toaster.create({
        title: "URL updated successfully",
        type: "success",
        duration: 3000,
      });
    },
    onError: (error: unknown) => {
      toaster.create({
        title: "Failed to update URL",
        description:
          (error as { response?: { data?: { message?: string } } })?.response
            ?.data?.message || (error as Error)?.message,
        type: "error",
        duration: 5000,
      });
    },
  });

  // Delete URL mutation
  const deleteURLMutation = useMutation({
    mutationFn: (id: number) => urlApi.deleteURL(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["urls"] });
      setSelectedURL(null);
      toaster.create({
        title: "URL deleted successfully",
        type: "success",
        duration: 3000,
      });
    },
    onError: (error: unknown) => {
      toaster.create({
        title: "Failed to delete URL",
        description:
          (error as { response?: { data?: { message?: string } } })?.response
            ?.data?.message || (error as Error)?.message,
        type: "error",
        duration: 5000,
      });
    },
  });

  // Handle create URL
  const handleCreateURL = () => {
    if (!createForm.original_url.trim()) {
      toaster.create({
        title: "URL is required",
        type: "error",
        duration: 3000,
      });
      return;
    }

    const data: CreateURLRequest = {
      original_url: createForm.original_url,
      short_code: createForm.short_code || undefined,
      expires_at: createForm.expires_at || undefined,
    };

    createURLMutation.mutate(data);
  };

  // Handle edit URL
  const handleEditURL = () => {
    if (!selectedURL) return;

    const data: UpdateURLRequest = {
      original_url: editForm.original_url || undefined,
      expires_at: editForm.expires_at || undefined,
      is_active: editForm.is_active,
    };

    updateURLMutation.mutate({ id: selectedURL.id, data });
  };

  // Handle delete URL
  const handleDeleteURL = (id: number) => {
    deleteURLMutation.mutate(id);
  };

  // Open edit form
  const openEditForm = (url: URL) => {
    setSelectedURL(url);
    setEditForm({
      original_url: url.original_url,
      expires_at: url.expires_at || "",
      is_active: url.is_active,
    });
    setShowEditForm(true);
  };

  // Copy to clipboard
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toaster.create({
      title: "Copied to clipboard",
      type: "success",
      duration: 2000,
    });
  };

  // Get short URL
  const getShortURL = (shortCode: string) => {
    const baseURL = import.meta.env.VITE_API_URL?.replace('/api/v1', '') || 'https://myapp-1757744589.fly.dev';
    return `${baseURL}/${shortCode}`;
  };

  if (error) {
    return (
      <Center p={8}>
        <Text color="red.500">
          Error loading URLs: {(error as Error).message}
        </Text>
      </Center>
    );
  }

  return (
    <Box p={6}>
      <VStack gap={6} align="stretch">
        {/* Header */}
        <HStack justify="space-between">
          <Heading size="lg">URL Manager</Heading>
          <Button
            colorPalette="blue"
            onClick={() => setShowCreateForm(!showCreateForm)}
          >
            {showCreateForm ? "Cancel" : "Create New URL"}
          </Button>
        </HStack>

        {/* Create Form */}
        {showCreateForm && (
          <Box p={4} bg="gray.50" borderRadius="md">
            <VStack gap={4}>
              <Heading size="md">Create New URL</Heading>
              <Box w="full">
                <Text mb={2}>Original URL *</Text>
                <Input
                  placeholder="https://example.com"
                  value={createForm.original_url}
                  onChange={(e) =>
                    setCreateForm({
                      ...createForm,
                      original_url: e.target.value,
                    })
                  }
                />
              </Box>
              <Box w="full">
                <Text mb={2}>Custom Short Code (optional)</Text>
                <Input
                  placeholder="abc123"
                  value={createForm.short_code}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, short_code: e.target.value })
                  }
                />
              </Box>
              <Box w="full">
                <Text mb={2}>Expiration Date (optional)</Text>
                <Input
                  type="datetime-local"
                  value={createForm.expires_at}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, expires_at: e.target.value })
                  }
                />
              </Box>
              <HStack>
                <Button
                  variant="outline"
                  onClick={() => setShowCreateForm(false)}
                >
                  Cancel
                </Button>
                <Button
                  colorPalette="blue"
                  onClick={handleCreateURL}
                  loading={createURLMutation.isPending}
                >
                  Create URL
                </Button>
              </HStack>
            </VStack>
          </Box>
        )}

        {/* Edit Form */}
        {showEditForm && selectedURL && (
          <Box p={4} bg="blue.50" borderRadius="md">
            <VStack gap={4}>
              <Heading size="md">Edit URL</Heading>
              <Box w="full">
                <Text mb={2}>Original URL</Text>
                <Input
                  placeholder="https://example.com"
                  value={editForm.original_url}
                  onChange={(e) =>
                    setEditForm({ ...editForm, original_url: e.target.value })
                  }
                />
              </Box>
              <Box w="full">
                <Text mb={2}>Expiration Date</Text>
                <Input
                  type="datetime-local"
                  value={editForm.expires_at}
                  onChange={(e) =>
                    setEditForm({ ...editForm, expires_at: e.target.value })
                  }
                />
              </Box>
              <HStack>
                <Button
                  variant="outline"
                  onClick={() => setShowEditForm(false)}
                >
                  Cancel
                </Button>
                <Button
                  colorPalette="blue"
                  onClick={handleEditURL}
                  loading={updateURLMutation.isPending}
                >
                  Update URL
                </Button>
              </HStack>
            </VStack>
          </Box>
        )}

        {/* Search */}
        <HStack>
          <Input
            placeholder="Search URLs..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            maxW="300px"
          />
        </HStack>

        {/* URLs List */}
        <Box>
          <Heading size="md" mb={4}>
            Your URLs
          </Heading>
          {isLoading ? (
            <Center p={8}>
              <Spinner size="lg" />
            </Center>
          ) : (
            <VStack gap={4} align="stretch">
              {urlsData?.data?.map((url) => (
                <Box
                  key={url.id}
                  p={4}
                  border="1px solid"
                  borderColor="gray.200"
                  borderRadius="md"
                  bg="white"
                >
                  <VStack gap={2} align="stretch">
                    <HStack justify="space-between">
                      <Text fontWeight="bold" fontSize="sm">
                        {url.original_url}
                      </Text>
                      <Badge colorPalette={url.is_active ? "green" : "red"}>
                        {url.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </HStack>

                    <HStack>
                      <Text fontFamily="mono" fontSize="sm">
                        {getShortURL(url.short_code)}
                      </Text>
                      <IconButton
                        size="xs"
                        variant="ghost"
                        onClick={() =>
                          copyToClipboard(getShortURL(url.short_code))
                        }
                      >
                        📋
                      </IconButton>
                    </HStack>

                    <HStack justify="space-between">
                      <Text fontSize="sm" color="gray.600">
                        Clicks: {url.click_count} | Created:{" "}
                        {new Date(url.created_at).toLocaleDateString()}
                      </Text>
                      <HStack>
                        <Button
                          size="xs"
                          variant="outline"
                          onClick={() => openEditForm(url)}
                        >
                          Edit
                        </Button>
                        <Button
                          size="xs"
                          variant="outline"
                          colorPalette="red"
                          onClick={() => handleDeleteURL(url.id)}
                          loading={deleteURLMutation.isPending}
                        >
                          Delete
                        </Button>
                      </HStack>
                    </HStack>
                  </VStack>
                </Box>
              ))}
            </VStack>
          )}

          {/* Pagination */}
          {urlsData?.pagination && (
            <HStack justify="center" mt={4}>
              <Button
                size="sm"
                variant="outline"
                onClick={() => setPage(page - 1)}
                disabled={page <= 1}
              >
                Previous
              </Button>
              <Text>
                Page {page} of {urlsData.pagination.total_pages}
              </Text>
              <Button
                size="sm"
                variant="outline"
                onClick={() => setPage(page + 1)}
                disabled={page >= urlsData.pagination.total_pages}
              >
                Next
              </Button>
            </HStack>
          )}
        </Box>
      </VStack>
    </Box>
  );
}
