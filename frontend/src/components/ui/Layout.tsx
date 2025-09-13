import { Box, Container } from "@chakra-ui/react";
import { Header } from "./Header";

export const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <Container minH="100vh" bg="gray.50">
      <Header />
      <Box maxW="1200px" mx="auto">
        {children}
      </Box>
    </Container>
  );
};
