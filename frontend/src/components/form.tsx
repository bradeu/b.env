"use client";

import { Input, Stack } from "@chakra-ui/react";
import { Field } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";

const schema = z.object({
  name: z.string().nonempty(),
  key: z.string().nonempty(),
});

type FormData = z.infer<typeof schema>;

export default function Form() {
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({ resolver: zodResolver(schema) });

  const onSubmit = async ({ name, key }: FormData) => {
    const validate = schema.safeParse({ name, key });
    if (validate.success) {
      // example of api call from client
      const res = await axios
        .post(`/api/send`, {
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
      <Stack gap="8" maxW="sm" css={{ "--field-label-width": "96px" }}>
        <Field
          orientation="horizontal"
          label="Name"
          invalid={Boolean(errors.name)}
        >
          <Input placeholder="John" flex="1" {...register("name")} />
        </Field>
        <Field
          orientation="horizontal"
          label="Key"
          invalid={Boolean(errors.key)}
        >
          <Input placeholder="Doe" flex="1" {...register("key")} />
        </Field>
        <Button type="submit" loading={isSubmitting}>
          Submit
        </Button>
      </Stack>
    </form>
  );
}
