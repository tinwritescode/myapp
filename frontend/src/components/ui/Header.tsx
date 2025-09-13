import { Button, Flex, HStack, Text } from "@chakra-ui/react";
import { Link, useLocation } from "react-router-dom";
import { useAuthStore } from "@/hooks/useAuthStore";

const menuList = [
  {
    label: "Shorten",
    path: "/shorten",
  },
  {
    label: "Manage",
    path: "/manage",
  },
];

export const Header = () => {
  const { isLoggedIn, user, logout } = useAuthStore();
  const location = useLocation();

  return (
    <Flex justifyContent="space-between" alignItems="center" p={4}>
      <Link to="/">
        <Text fontSize="xl" fontWeight="semibold">
          URL Shortener
        </Text>
      </Link>

      <HStack gap={4}>
        {isLoggedIn ? (
          <>
            {menuList.map((menu) => (
              <Link key={menu.path} to={menu.path}>
                <Button
                  variant={location.pathname === menu.path ? "solid" : "ghost"}
                  colorPalette={
                    location.pathname === menu.path ? "blue" : "gray"
                  }
                  size="sm"
                >
                  {menu.label}
                </Button>
              </Link>
            ))}
            <Text fontSize="sm">Welcome, {user?.email}</Text>
            <Button colorPalette="red" size="sm" onClick={logout}>
              Logout
            </Button>
          </>
        ) : (
          <Text fontSize="sm">Please login to access features</Text>
        )}
      </HStack>
    </Flex>
  );
};
