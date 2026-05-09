package com.bryan.system.mapper;

import org.junit.jupiter.api.Test;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import javax.xml.parsers.DocumentBuilderFactory;
import java.io.InputStream;
import java.lang.reflect.Method;
import java.util.HashSet;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class MapperXmlBindingTest {

    @Test
    void shouldExposeMappedStatementsForAllMapperMethods() throws Exception {
        assertMapperXmlBindings(UserMapper.class, "mapper/UserMapper.xml");
        assertMapperXmlBindings(UserProfileMapper.class, "mapper/UserProfileMapper.xml");
        assertMapperXmlBindings(UserRoleMapper.class, "mapper/UserRoleMapper.xml");
    }

    private void assertMapperXmlBindings(Class<?> mapperType, String xmlPath) throws Exception {
        Set<String> statementIds = parseStatementIds(xmlPath, mapperType.getName());
        for (Method method : mapperType.getDeclaredMethods()) {
            assertTrue(statementIds.contains(method.getName()),
                    "Missing XML statement for " + mapperType.getName() + "." + method.getName());
        }
    }

    private Set<String> parseStatementIds(String xmlPath, String expectedNamespace) throws Exception {
        try (InputStream in = Thread.currentThread().getContextClassLoader().getResourceAsStream(xmlPath)) {
            assertTrue(in != null, "XML resource not found: " + xmlPath);

            DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
            factory.setNamespaceAware(false);
            Document document = factory.newDocumentBuilder().parse(in);

            Element mapper = document.getDocumentElement();
            assertEquals("mapper", mapper.getTagName(), "Root element must be mapper: " + xmlPath);
            assertEquals(expectedNamespace, mapper.getAttribute("namespace"),
                    "Mapper namespace mismatch: " + xmlPath);

            Set<String> ids = new HashSet<>();
            NodeList children = mapper.getChildNodes();
            for (int i = 0; i < children.getLength(); i++) {
                Node node = children.item(i);
                if (node.getNodeType() != Node.ELEMENT_NODE) {
                    continue;
                }
                String tag = node.getNodeName();
                if (!"select".equals(tag) && !"insert".equals(tag) && !"update".equals(tag) && !"delete".equals(tag)) {
                    continue;
                }
                Element element = (Element) node;
                String id = element.getAttribute("id");
                if (id != null && !id.isBlank()) {
                    ids.add(id);
                }
            }
            return ids;
        }
    }
}
