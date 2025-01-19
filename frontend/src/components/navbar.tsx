import { Center, Heading } from "@chakra-ui/react";

export default function Navbar() {
  return (
    <Center position={"fixed"} top={0} width={"100%"} height={"2.5rem"}>
      <Heading fontSize={"2xl"}>
        {"<"}b.env{">"}
      </Heading>
    </Center>
  );
}
