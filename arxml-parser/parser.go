package arxml

import (
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"os"
	"strconv"
	"strings"
)

var (
	NoDataError = errors.New("No data")
)

func parseFile(file *os.File) (*xmlquery.Node, error) {
	doc, err := xmlquery.Parse(file)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func parseXml(filePath string) (*xmlquery.Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return parseFile(file)
}

func getName(node *xmlquery.Node) string {
	lst := xmlquery.Find(node, "/SHORT-NAME")
	if lst != nil && len(lst) == 0 {
		return ""
	}
	return lst[0].FirstChild.Data
}

func listPackages(root *xmlquery.Node) [] *xmlquery.Node {
	if root == nil { return nil }
	return xmlquery.Find(root, "//AR-PACKAGES/AR-PACKAGE")
}

func getPackage(root *xmlquery.Node, name string) *xmlquery.Node {
	if root == nil { return nil }
	for _, p := range listPackages(root) {
		if getName(p) == name {
			return p
		}
	}
	return nil
}

func getItem(nodes []*xmlquery.Node, name string) *xmlquery.Node {
	for _, node := range nodes {
		if getName(node) == name {
			return node
		}
	}
	return nil
}

func getObjects(node *xmlquery.Node, name string) []*xmlquery.Node {
	if node == nil { return nil }
	return xmlquery.Find(node, fmt.Sprintf("/%s", name))
}

func getObjectsInside(node *xmlquery.Node, name string) []*xmlquery.Node {
	if node == nil { return nil }
	return xmlquery.Find(node, fmt.Sprintf("//%s", name))
}

func getText(node *xmlquery.Node) (string, error) {
	if node == nil { return "", NoDataError }
	return strings.TrimSpace(node.FirstChild.Data), nil
}

func getHeadText(nodes []*xmlquery.Node) (string, error) {
	if nodes == nil || len(nodes) ==0 { return "", NoDataError }
	return strings.TrimSpace(nodes[0].FirstChild.Data), nil
}

func getIntText(str string, err error) int32 {
	var ret int32 = 0
	if err == nil {
		val, _ := strconv.ParseInt(str, 10, 32)
		ret = int32(val)
	}
	return ret
}

func getIntValue(str string) (int32, error) {
	val, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

func getFloatText(str string, err error) float32 {
	var ret float32 = 0.0
	if err == nil {
		val, _ := strconv.ParseFloat(str, 32)
		ret = float32(val)
	}
	return ret
}

func getFloatValue(str string) (float32, error) {
	val, err := strconv.ParseFloat(str, 32)
	if err != nil { return 0.0, err }
	return float32(val), nil
}


func getNetwork(root *xmlquery.Node) []Network {
	if root == nil { return nil }
	clusters := getPackage(getPackage(root,"Topology"), "Clusters")
	ethernets := getObjects(getObjects(clusters, "ELEMENTS")[0],"ETHERNET-CLUSTER")
	ethernet := getItem(ethernets, "Ethernet_Cluster")
	channels := getObjectsInside(ethernet, "ETHERNET-PHYSICAL-CHANNEL")
	networks := make([]Network,0)
	for _, ch := range channels {
		name := getName(ch)
		vid := getIntText(getHeadText(xmlquery.Find(ch, "/VLAN/VLAN-IDENTIFIER")))
		pduRef := make(map[string]int32)
		identifiers := xmlquery.Find(ch, "//SOCKET-CONNECTION-IPDU-IDENTIFIER")

		for _, node := range identifiers {
			idStr, _ := getHeadText(xmlquery.Find(node, "/HEADER-ID"))
			ref, _ := getHeadText(xmlquery.Find(node, "/PDU-TRIGGERING-REF"))
			if len(idStr) > 0 && len(ref) > 0 {
				id, err := getIntValue(idStr)
				if err == nil {
					refName := strings.Split(ref, "/")
					pduRef[refName[len(refName)-1]] = id
				}
			}
		}
		pdus := make([]PduRef, 0)
		triggers := xmlquery.Find(ch, "//PDU-TRIGGERING")
		for _, node := range triggers {
			pname := getName(node)
			ref, err := getHeadText(xmlquery.Find(node, "/I-PDU-REF"))
			if err == nil && len(pname) >0 {
				id := pduRef[pname]
				pdus = append(pdus, newPduRef(pname, ref, id))
			}
		}
		networks = append(networks, newNetwork(name, vid, pdus))
	}
	return networks
}

func getISignal(root *xmlquery.Node) []ISignal {
	isignals := make([]ISignal, 0)
	if root == nil { return nil }
	signals := getPackage(getPackage(root, "Communication"), "Signals")
	sigs := getObjectsInside(signals, "I-SIGNAL")
	for _, sig := range sigs {
		name := getName(sig)
		desc, _ := getHeadText(xmlquery.Find(sig, "/DESC"))
		length := getIntText(getHeadText(xmlquery.Find(sig,"/LENGTH")))
		value  := getFloatText(getHeadText(xmlquery.Find(sig,  "//VALUE")))
		ref, _ := getHeadText(xmlquery.Find(sig, "//COMPU-METHOD-REF"))
		typeRef, err := getHeadText(xmlquery.Find(sig,  "//BASE-TYPE-REF"))
		var signed bool = false
		if err == nil && strings.Contains(typeRef, "SINT") {
			signed = true
		}
		isignals = append(isignals, NewISignal(name, length, desc, ref, value, signed))
	}
	return isignals
}

func getDataTypes(root *xmlquery.Node) []ComputeMethod {
	computeMethods := make([]ComputeMethod, 0)
	compus := xmlquery.Find(getPackage(getPackage(root, "DataTypes"), "CompuMethods"),  "//COMPU-METHOD")
	for _, compu := range compus {
		name := getName(compu)
		category, caterr := getHeadText(xmlquery.Find(compu, "/CATEGORY"))
		ref, referr := getHeadText(xmlquery.Find(compu, "/UNIT-REF"))
		if caterr == nil && category != "IDENTICAL" {
			var unit string = ""
			if referr == nil {
				unit = ref[len("/DataTypes/Units/"):]
			}
			compuScale := make([]CompuScale, 0)

			for _, n := range xmlquery.Find(compu, "/COMPU-INTERNAL-TO-PHYS") {
				for _, scale := range xmlquery.Find(n, "//COMPU-SCALE") {
					label, err := getHeadText(xmlquery.Find(scale, "/SHORT-LABEL"))
					if err == nil {
						min := getFloatText(getHeadText(xmlquery.Find(scale,  "/LOWER-LIMIT")))
						max := getFloatText((getHeadText(xmlquery.Find(scale,  "/UPPER-LIMIT"))))
						nums := make([]float32, 0)
						for _, vn := range xmlquery.Find(scale, "//COMPU-NUMERATOR/V") {
							num := getFloatText(getText(vn))
							nums = append(nums, num)
						}
						denominator := getFloatText(getHeadText(xmlquery.Find(scale,  "//COMPU-DENOMINATOR/V")))
						constant, _ := getHeadText(xmlquery.Find(scale,  "//VT"))
						compuScale = append(compuScale, NewCompuScale(label, min, max, NewCompuNum(nums[0], nums[1]), denominator, constant))
					}
				}
			}
			computeMethods = append(computeMethods, NewComputeMethod(name,category,unit, compuScale))
		}
	}
	return computeMethods
}
