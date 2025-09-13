"use client";

import {
  Toaster as ChakraToaster,
  Portal,
  Spinner,
  Stack,
  Toast,
  createToaster,
} from "@chakra-ui/react";

import { type CreateToasterReturn } from "@ark-ui/react";

export const toaster = createToaster({
  placement: "bottom-end",
  pauseOnPageIdle: true,
}) as unknown as CreateToasterReturn;

export const Toaster = () => {
  return (
    <Portal>
      {/* @ts-ignore */}
      <ChakraToaster toaster={toaster} insetInline={{ mdDown: "4" }}>
        {/* @ts-ignore */}
        {(toast: any) => (
          <Toast.Root width={{ md: "sm" }}>
            {toast.type === "loading" ? (
              <Spinner size="sm" color="blue.solid" />
            ) : (
              <Toast.Indicator />
            )}
            <Stack gap="1" flex="1" maxWidth="100%">
              {/* @ts-ignore */}
              {toast.title && <Toast.Title>{toast.title}</Toast.Title>}
              {/* @ts-ignore */}
              {toast.description && (
                <>
                  {/* @ts-ignore */}
                  <Toast.Description>{toast.description}</Toast.Description>
                </>
              )}
            </Stack>
            {/* @ts-ignore */}
            {toast.action && (
              <>
                {/* @ts-ignore */}
                <Toast.ActionTrigger>{toast.action.label}</Toast.ActionTrigger>
              </>
            )}
            {toast.closable && <Toast.CloseTrigger />}
          </Toast.Root>
        )}
      </ChakraToaster>
    </Portal>
  );
};
