"use client";

import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import { Input, Stack } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import { useForm } from "react-hook-form";
import z from "zod";

const schema = z.object({
  address: z.string().min(1),
  name: z.string().min(1),
  key: z.string().min(1),
});

type FormData = z.infer<typeof schema>;

export default function Form({ address }: { address: string }) {
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({ resolver: zodResolver(schema) });

  const onSubmit = async ({ name, key }: FormData) => {
    const validate = schema.safeParse({ address, name, key });
    if (validate.success) {
      // example of api call from client
      const res = await axios
        .post(`/api/send`, {
          address,
          name,
          key,
        })
        .catch((e) => {
          if (axios.isAxiosError(e)) {
            console.log("HTTPS req is not working");
          }
        });
      console.log(res);
    }
  };
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Stack
        gap="8"
        width="100%"
        height={"fit-content"}
        css={{ "--field-label-width": "96px" }}
      >
        <Field
          orientation="vertical"
          label="Name"
          invalid={Boolean(errors.name)}
          errorText={errors.name?.message}
        >
          <Input placeholder="John" flex="1" p={3} {...register("name")} />
        </Field>
        <Field
          orientation="vertical"
          label="Key"
          invalid={Boolean(errors.key)}
          errorText={errors.key?.message}
        >
          <Input placeholder="Doe" flex="1" p={3} {...register("key")} />
        </Field>
        <Button type="submit" loading={isSubmitting}>
          Submit
        </Button>
      </Stack>
    </form>
  );
}
