package links_test
//
//import (
//	. "github.com/jjmschofield/GoCrawl/internal/app/links"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"net/url"
//)
//
//var _ = XDescribe("Group", func() {
//	Describe("ToLinkGroup", func(){
//		var (
//			underTest func(links []Link) (group LinkGroup)
//		)
//
//		BeforeEach(func(){
//			underTest = ToLinkGroup
//		})
//
//		// This test is really quite lazy in testing a big switch statement in one go - but the code is quite simple and robust
//		// If you experience problems with this test (or you find it's got really big), boy scout to test each case individually
//		It("should parse links to a group structure", func(){
//			// Arrange
//			links := createLinksFake()
//
//			// Act
//			result := underTest(links)
//
//			// Assert
//			Expect(result.Internal[0]).To(Equal(links[0]))
//			Expect(result.InternalFile[0]).To(Equal(links[1]))
//			Expect(result.External[0]).To(Equal(links[2]))
//			Expect(result.ExternalFile[0]).To(Equal(links[3]))
//			Expect(result.Tel[0]).To(Equal(links[4]))
//			Expect(result.Mailto[0]).To(Equal(links[5]))
//			Expect(result.Unknown[0]).To(Equal(links[6]))
//		})
//	})
//
//	Describe("MarshalJson", func(){
//		It("should marshal links into ids", func(){
//			// Arrange
//			links := createLinksFake()
//			group := ToLinkGroup(links)
//			expected := []byte(`{"internal":["2d84fbc4b68dc0448c93a7d4deda20dc"],"internalFiles":["d9963df610562025fcbfd08269c726ba"],"external":["9b534e5fcc603649270bf40a4871c201"],"externalFiles":["9eba24981f0d28aa9e14c4791fa84202"],"tel":["eef84ec73f7ea7e38cb2cd2301c0900f"],"mailto":["b6adc74f00ae9b8ca5f274d6e57e672b"],"unknown":["52830e41c1915a8f4e53771e984eff90"]}`)
//
//			// Act
//			result, _ := group.MarshalJSON()
//
//			// Assert
//			Expect(result).To(Equal(expected))
//		})
//	})
//})
//
//func createLinksFake() (links []Link){
//	srcUrl, _ := url.Parse("https://www.google.co.uk/")
//
//	hrefs := []string{
//		"https://www.google.co.uk/internal/page",
//		"https://www.google.co.uk/internal/file.zip",
//		"https://www.bbc.co.uk/external/page",
//		"https://www.bbc.co.uk/external/file.pdf",
//		"tel:+440123456789",
//		"mailto:someone@somewhere.com",
//	}
//
//	links = FromHrefs(*srcUrl, hrefs)
//
//	unknownLink, _ := FromHref(*srcUrl, "https://some.wild/unkown/link/type")
//	unknownLink.Type = UnknownType
//	links = append(links, unknownLink)
//
//	return links
//}
