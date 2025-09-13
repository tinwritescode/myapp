import { ChakraProvider, createSystem, defaultConfig } from "@chakra-ui/react";
import { ThemeProvider } from "next-themes";

const system = createSystem(defaultConfig, {
  theme: {
    tokens: {
      fonts: {
        body: { value: "Poppins, sans-serif" },
        heading: { value: "Poppins, sans-serif" },
      },
    },
  },
});

export const Provider = ({ children }: { children: React.ReactNode }) => {
  return (
    <ChakraProvider value={system}>
      <ThemeProvider attribute="class">{children}</ThemeProvider>
    </ChakraProvider>
  );
};
