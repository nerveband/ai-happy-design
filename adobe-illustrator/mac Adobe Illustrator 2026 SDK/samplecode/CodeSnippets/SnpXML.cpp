//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 1987 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h" 
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"

/* XML document for the selected document. */
AIXMLNodeRef element = NULL;

/**
	Provides functions to create, modify, inspect, and delete XML nodes associated with 
	a document.
	Also provides functions to iterate through the metadata of the XML document of the selected document, 
	and also to delete them.
*/
class SnpXML
{
	public:		

		/* Metadata for the XML document of the selected document. */
		AIXMLNodeRef metadata ;

		/**
			Constructor for the snippet.
			It retrieves the XML document and the metadata element of the selected document.
			Creates if none exists.
		*/
		SnpXML();
		
		/**
			Demonstrates how to create various types of nodes.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateXMLNodes(void);		

		/**
			Creates the named node map for the metadata element.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateXMLMetadata(void);
		
		/**
			Modifies the XML node for the current document through various APIs.
			Also illustrates how to append and insert new nodes in an already existing nodelist.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ModifyXMLNode(void);

		/**
			Modifies the XML metadata for a node.
			Illustrates how to inspect and set new items in a named node map.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ModifyXMLMetadata(void);

		/**
			Displays the node being referred to by any member function, or otherwise.
			@param node IN which is to be displayed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DisplayXMLNode(const AIXMLNodeRef node);

		/**
			Displays the metadata associated with the XML document of the selected document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DisplayXMLMetadata(void);

		/**
			Removes the XML Nodes of the selected document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveXMLNodes(void);

		/**
			Removes the XML Named Nodes of the selected document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveXMLMetadata(void);		

		
		//------------------------------- Utility Functions -------------------------------

		/**
			Creates an element node.
			@return AIXMLNodeRef object as the newly created node.
		*/
		AIXMLNodeRef CreateXMLElement(void);

		/**
			Inspects the NodeList passed from other functions, and checks how many nodes it has and their related properties.
			This function hsa been called whenever a nodelist is to be inspected.
			@param  nodes IN  whose properties are to be inspected.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectXMLNodeList(const AIXMLNodeListRef nodes);


		/**
			Inspects the Node passed from other functions, and checks if it contains children and their related properties.
			This function has been called whenever a node is to be inspected.
			@param node IN whose properties are to be inspected.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectXMLNode(const AIXMLNodeRef node);

		/**
			Sets the specified attributes of the element node passed to the function.
			@param metadata IN/OUT node whose attributes are being set.
			@param xmlName IN XML name of the named node.
			@param xmlString IN XML value of the named node.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetXMLElementAttribute(AIXMLNodeRef metadata,const char* xmlName,const char* xmlString);

		/**
			Removes the specified attributes of the element node passed to the function.
			@param metadata IN/OUT node whose attributes are being removed.
			@param xmlName IN XML name of the named node.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveXMLElementAttribute(AIXMLNodeRef metadata,const char* xmlName);
		
};

/*
*/
SnpXML :: SnpXML(void)
{
	ASErr result = kNoErr;
	try {
		//Create the XML document for the current document.
		result = sAIXMLDocument->GetDocumentElement(&element);
		aisdk::check_ai_error(result);

		//Create the metadata element for the XML document of the current document.
		result = sAIXMLDocument->GetDocumentMetadata(&metadata);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
}

/*
*/
ASErr SnpXML::CreateXMLNodes(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		log->Write("");
		ASInt32 i = 1;
		log->WritePrintf("Node #%d",i);
		result = this->InspectXMLNode(element);
		aisdk::check_ai_error(result);

		log->Write("");
		log->WritePrintf("Node #%d",++i);
		result = this->DisplayXMLNode(metadata);
		aisdk::check_ai_error(result);

		//Create a new metadata text node.
		const char* string = "XMLTextNode";
		AIXMLNodeRef text;
		result = sAIXMLDocument->CreateTextNode(sAIEntry->FromString(string),&text);
		aisdk::check_ai_error(result);		
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		result = this->DisplayXMLNode(text);
		aisdk::check_ai_error(result);

		//Create a new metadata comment node.
		AIXMLNodeRef comment;
		const char* commentString = "Comment String";
		result = sAIXMLDocument->CreateComment(sAIEntry->FromString(commentString),&comment);
		aisdk::check_ai_error(result);	
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		result = this->DisplayXMLNode(comment);
		aisdk::check_ai_error(result);

		//Create a new metadata element.
		result = sAIXMLDocument->CreateElement(sAIXMLNode->NameFromUnicodeString(ai::UnicodeString("XML Element")),&element);
		aisdk::check_ai_error(result);
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		result = this->DisplayXMLNode(element);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::CreateXMLMetadata(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		//Set the String attribute for named node.
		result = this->SetXMLElementAttribute(metadata,"#String","String value for Named Node");
		aisdk::check_ai_error(result);
		log->Write("");
		ASInt32 i = 1;
		log->WritePrintf("Node #%d",i);
		log->Write("Creating metadata attribute node: ");
		result = this->DisplayXMLNode(metadata);
		aisdk::check_ai_error(result);	

		//Set the Integer attribute for named node.
		result = this->SetXMLElementAttribute(metadata,"#Integer","Integer value for Named Node");
		aisdk::check_ai_error(result);
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		log->Write("Creating metadata attribute node: ");
		result = this->DisplayXMLNode(metadata);
		aisdk::check_ai_error(result);	


		//Set the UnicodeString attribute for named node.
		result = this->SetXMLElementAttribute(metadata,"#UnicodeString","UnicodeString value for Named Node");
		aisdk::check_ai_error(result);
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		log->Write("Creating metadata attribute node: ");
		result = this->DisplayXMLNode(metadata);
		aisdk::check_ai_error(result);	


		//Set the Real attribute for named node.
		result = this->SetXMLElementAttribute(metadata,"#Real","Real value for Named Node");
		aisdk::check_ai_error(result);
		log->Write("");
		log->WritePrintf("Node #%d",++i);
		log->Write("Creating metadata attribute node: ");
		result = this->DisplayXMLNode(metadata);
		aisdk::check_ai_error(result);	

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::ModifyXMLNode(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	
		//Get Child Nodes.
		AIXMLNodeListRef elements = NULL;
		result = sAIXMLNode->GetChildNodes(element, &elements);
		aisdk::check_ai_error(result);

		//Get the number of child nodes.
		ai::int32 length = 0;
		result = sAIXMLNodeList->GetLength(elements, &length);
		aisdk::check_ai_error(result);
		log->WritePrintf("Number of nodes originally: %d",length);

		result = this->InspectXMLNodeList(elements);
		aisdk::check_ai_error(result);

		//Get Node Type.
		AIXMLNodeType type = kAIXMLUnknownNode;
		result = sAIXMLNode->GetNodeType(element, &type);
		aisdk::check_ai_error(result);

		if(type == kAIXMLDocumentNode){

			AIXMLNodeRef newNode = this->CreateXMLElement();
			aisdk::check_ai_error(result);
			
			//Append a new Element child node to the node list.
			result = sAIXMLNode->AppendChild(element,newNode);
			aisdk::check_ai_error(result);

			log->Write("Child has been appended");
			log->Write("Inspecting parent node again after Appending");
			result = this->InspectXMLNode(element);
			aisdk::check_ai_error(result);

			//Illustration of inserting a new child in an XML node.
			AIXMLNodeRef newchild = this->CreateXMLElement();
			result = sAIXMLNode->InsertBefore(element,newchild,newNode);
			aisdk::check_ai_error(result);
			log->Write("Child has been inserted");
			log->Write("Inspecting parent node again after Inserting");
			result = this->InspectXMLNode(element);
			aisdk::check_ai_error(result);
			
		}
		//Make an exact copy of the parent node.
		AIXMLNodeRef clone;
		result = sAIXMLNode->Clone(element,&clone);
		aisdk::check_ai_error(result);

		//Inspect the clone node.
		result = sAIXMLNode->GetChildNodes(clone, &elements);
		aisdk::check_ai_error(result);
		result = sAIXMLNodeList->GetLength(elements, &length);
		aisdk::check_ai_error(result);

		log->Write("");
		log->WritePrintf("Number of nodes in Clone node : %d",length);
		result = this->InspectXMLNode(clone);
		aisdk::check_ai_error(result);

		//Creates a new text node.
		const char* string = "XMLTextNode";
		AIXMLNodeRef text;
		result = sAIXMLDocument->CreateTextNode(sAIEntry->FromString(string),&text);
		aisdk::check_ai_error(result);	
		log->Write("");
		log->Write("Creating Text Node");
		result = this->InspectXMLNode(text);
		aisdk::check_ai_error(result);

		//Create a new Comment node.
		AIXMLNodeRef comment;
		const char* commentString = "Comment String in Clone Node";
		result = sAIXMLDocument->CreateComment(sAIEntry->FromString(commentString),&comment);
		aisdk::check_ai_error(result);
		log->Write("");
		log->Write("Creating Comment Node");
		result = this->InspectXMLNode(comment);
		aisdk::check_ai_error(result);

		//Compare the text and comment nodes for equality.
		ai::int32 equal;
		result = sAIXMLNode->Compare(text,comment,true,&equal);
		aisdk::check_ai_error(result);
		log->Write("");
		if(equal > 0)
			log->Write("Text and Comment nodes being compared are equal");
		else
			log->Write("Text and Comment nodes being compared are not equal");		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::ModifyXMLMetadata(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		//Retrieve a named node map containing the attributes associated with the XML node. 
		AIXMLNamedNodeMapRef xmlAttributes;
		result = sAIXMLNode->GetAttributes(metadata,&xmlAttributes);
		aisdk::check_ai_error(result);

		//Get single nodes from the named node map.
		AIXMLNodeRef node;
		result = sAIXMLNamedNodeMap->GetNamedItem(xmlAttributes,sAIXMLNode->NameFromString("#UnicodeString"),&node);
		aisdk::check_ai_error(result);

		log->Write("");
		log->Write("Retrieved Named Node :");
		result = this->DisplayXMLNode(node);
		aisdk::check_ai_error(result);	

		//Create a new attribute node to insert in named node map.
		const char* xmlName = "#New Attribute Node";
		const char* xmlString = "New Attribute Value";
		AIXMLNodeRef attribute;
		result = sAIXMLDocument->CreateAttribute(sAIXMLNode->NameFromString(xmlName),sAIEntry->FromString(xmlString),&attribute);
		aisdk::check_ai_error(result);	

		//Set the created item in the named node map.
		result = sAIXMLNamedNodeMap->SetNamedItem(xmlAttributes,attribute);
		aisdk::check_ai_error(result);

		log->Write("");
		log->Write("Added named node :");
		result = this->DisplayXMLNode(attribute);
		aisdk::check_ai_error(result);	

		log->Write("");
		log->Write("Inspecting named node list after addition of one more named node:");
		result = this->DisplayXMLMetadata();
		aisdk::check_ai_error(result);	
		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpXML :: DisplayXMLNode(AIXMLNodeRef node)
{
	ASErr result = kNoErr;
	// Write output to the log.
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	SnippetRunnerLog::Indent indent;
	try {

		//Get node type.
		AIXMLNodeType type = kAIXMLUnknownNode;
		result = sAIXMLNode->GetNodeType(node, &type);
		aisdk::check_ai_error(result);				

		switch(type){
			case kAIXMLElementNode:
				log->Write("Node Type : Element Node");
				break;
			case kAIXMLAttributeNode:
				log->Write("Node Type : Attribute Node");
				break;
			case kAIXMLTextNode:
				log->Write("Node Type : Text Node");
				break;
			case kAIXMLCDATASectionNode:
				log->Write("Node Type : CDATASection Node");
				break;
			case kAIXMLCommentNode:
				log->Write("Node Type : Comment Node");
				break;
			case kAIXMLDocumentNode:
				log->Write("Node Type :  Document Node");
				break;
			default:
				log->Write("Node Type : Unknown ");
				break;
			}
			//Get node name.
			AIXMLName name;
			result = sAIXMLNode->GetNodeName(node,&name);
			aisdk::check_ai_error(result);
			
			ai::UnicodeString nodeName;
			result = sAIXMLNode->UnicodeStringFromName(name,nodeName);
			aisdk::check_ai_error(result);
			log->Write(ai::UnicodeString("Node Name : ").append(nodeName));

			//Get node value.
			AIXMLString value;
			result = sAIXMLNode->GetNodeValue(node,&value);
			aisdk::check_ai_error(result);

			const char *stringValue;
			result = sAIEntry->ToString(value,&stringValue);
			aisdk::check_ai_error(result);
			if(*stringValue)
				log->WritePrintf("Node Value : %s ",stringValue);
			else
				log->Write("Node Value : NULL");
			
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::DisplayXMLMetadata(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		//Retrieve a named node map containing the attributes associated with the XML node. 
		AIXMLNamedNodeMapRef xmlAttributes;
		result = sAIXMLNode->GetAttributes(metadata,&xmlAttributes);
		aisdk::check_ai_error(result);

		//Retrieve the size of the named node map.
		ai::int32 length;
		result = sAIXMLNamedNodeMap->GetLength(xmlAttributes,&length);
		aisdk::check_ai_error(result);
		log->WritePrintf("The length of the metadata named node list is : %d",length);

		for(ai::int32 i = 0; i < length; i++){
			//Get single nodes from the named node map.
			AIXMLNodeRef node;
			result = sAIXMLNamedNodeMap->GetItem(xmlAttributes,i,&node);
			aisdk::check_ai_error(result);
	
			log->Write("");
			log->WritePrintf("Node #%d",i);
			result = this->DisplayXMLNode(node);
			aisdk::check_ai_error(result);			
		}	
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::RemoveXMLNodes(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	
		//Get Child Nodes.
		AIXMLNodeListRef elements = NULL;
		result = sAIXMLNode->GetChildNodes(element, &elements);
		aisdk::check_ai_error(result);

		//Get the number of child nodes.
		ai::int32 length = 0;
		result = sAIXMLNodeList->GetLength(elements, &length);
		aisdk::check_ai_error(result);
		log->WritePrintf("Number of nodes before Releasing : %d",length);

		for(ai::int32 i = length-1; i >= 0; i--)
		{
			//Get each child node.
			AIXMLNodeRef child;
			result = sAIXMLNodeList->GetItem(elements, i, &child);
			aisdk::check_ai_error(result);
			
			log->WritePrintf("Removing node #%d",i);
			result = this->DisplayXMLNode(child);
			aisdk::check_ai_error(result);

			//Remove the child node.
			result = sAIXMLNode->RemoveChild(element,child);
			aisdk::check_ai_error(result);			
		}
		log->Write("");
		log->Write("After removing the Nodes");
		result = this->InspectXMLNodeList(elements);
		aisdk::check_ai_error(result);

		//Release the memory associated with the node element.
		ai::int32 ref = sAIXMLNode->Release(element);
		if(ref)
			log->Write("Node has been released");
		else
			log->Write("Node has not been released");
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::RemoveXMLMetadata(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		//Retrieve a named node map containing the attributes associated with the XML node. 
		AIXMLNamedNodeMapRef xmlAttributes;
		result = sAIXMLNode->GetAttributes(metadata,&xmlAttributes);
		aisdk::check_ai_error(result);

		//Retrieve the size of the named node map.
		ai::int32 length;
		result = sAIXMLNamedNodeMap->GetLength(xmlAttributes,&length);
		aisdk::check_ai_error(result);
		log->WritePrintf("The length of the metadata named node list is : %d",length);

		for(ai::int32 i = length-1; i >= 0; i--)
		{
			//Get single nodes from the named node map.
			AIXMLNodeRef node;
			result = sAIXMLNamedNodeMap->GetItem(xmlAttributes,i,&node);
			aisdk::check_ai_error(result);

			log->Write("");
			log->WritePrintf("Removing named node #%d",i);
			result = this->DisplayXMLNode(node);
			aisdk::check_ai_error(result);

			AIXMLName xmlName;
			result = sAIXMLNode->GetNodeName(node,&xmlName);
			aisdk::check_ai_error(result);

			//Remove the named node.
			result = this->RemoveXMLElementAttribute(metadata,sAIXMLNode->StringFromName(xmlName));
			aisdk::check_ai_error(result);
		}		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
AIXMLNodeRef SnpXML::CreateXMLElement(void)
{
	ASErr result = kNoErr;
	AIXMLNodeRef newelement;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		//Create the XML document for the current document.
		ai::UnicodeString newName = SnippetRunnerParameter::Instance()->GetString(ai::UnicodeString("Enter new name for the New Node"),ai::UnicodeString("New Node"));
		result = sAIXMLDocument->CreateElement(sAIXMLNode->NameFromUnicodeString(newName),&newelement);
		aisdk::check_ai_error(result);

		log->Write("");
		log->Write("Created new Node");
		result = this->InspectXMLNode(newelement);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return newelement;
}

/*
*/
ASErr SnpXML :: InspectXMLNodeList(AIXMLNodeListRef nodes)
{
	ASErr result = kNoErr;
	try {
			//Get the number of nodes in the NodeList.
			ai::int32 length = 0;
			result = sAIXMLNodeList->GetLength(nodes, &length);
			aisdk::check_ai_error(result);

			// Write output to the log.
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->WritePrintf("Number of nodes in NodeList: %d",length);

			for( ai::int32 i = 0; i < length; i++)
			{
				//Retrieve the nodes in the nodelist.
				AIXMLNodeRef child;
				result = sAIXMLNodeList->GetItem(nodes, i, &child);
				aisdk::check_ai_error(result);

				//Inspect the individual node.
				result = this->InspectXMLNode(child);
				aisdk::check_ai_error(result);
			}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML :: InspectXMLNode(AIXMLNodeRef node)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;

		//Get the type of node.
		AIXMLNodeType type = kAIXMLUnknownNode;
		result = sAIXMLNode->GetNodeType(node, &type);
		aisdk::check_ai_error(result);

		//Display the node details.
		result = this->DisplayXMLNode(node);
		aisdk::check_ai_error(result);

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		//Check if the node has children.
		AIBoolean haschildren;
		result = sAIXMLNode->HasChildNodes(node,&haschildren);
		aisdk::check_ai_error(result);		

		if(haschildren){
			log->Write("This node has children");
			//Get Child nodes of the current node.
			AIXMLNodeListRef elements;
			result = sAIXMLNode->GetChildNodes(node, &elements);
			aisdk::check_ai_error(result);

			//Get the number of children.
			ai::int32 length = 0;
			result = sAIXMLNodeList->GetLength(elements, &length);
			aisdk::check_ai_error(result);
			log->WritePrintf("Number of nodes in the nodelist being inspected : %d",length);

			for( ai::int32 i = 0; i < length; i++)
			{
				//Get each child node.
				AIXMLNodeRef child;
				result = sAIXMLNodeList->GetItem(elements, i, &child);
				aisdk::check_ai_error(result);

				//Display the node details.
				log->WritePrintf("Child Node #%d: ",i);
				result = this->DisplayXMLNode(child);
				aisdk::check_ai_error(result);				
			}
		}
		else{
			log->Write("This node has no children");

		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::SetXMLElementAttribute(AIXMLNodeRef metadata,const char* xmlName,const char* xmlString)
{
	ASErr result = kNoErr;
	try {
		//Set the value of the XML's node's named attribute.
		result = sAIXMLElement->SetAttribute(metadata,sAIXMLNode->NameFromString(xmlName),sAIEntry->FromString(xmlString));
		aisdk::check_ai_error(result);	
			
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpXML::RemoveXMLElementAttribute(AIXMLNodeRef metadata,const char* xmlName)
{
	ASErr result = kNoErr;
	try {
		//Remove the named attribute associated with the XML node.
		result = sAIXMLElement->RemoveAttribute(metadata,sAIXMLNode->NameFromString(xmlName));
		aisdk::check_ai_error(result);	
			
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// ---------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpXML available to the SnippetRunner framework.
 */
class _SnpRunnableXML : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableXML () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableXML () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation - must
			be one of the operation names returned by GetOperations.
		*/
		std::string GetDefaultOperationName() const;

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/* Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableXML::GetName() const
{
	return "XML";
}

/*
*/
std::string _SnpRunnableXML::GetDescription() const 
{
	return "Demonstrates XML Operations";
}

/*
*/
SnpRunnable::Operations _SnpRunnableXML::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateXMLNodes", "document", kSnpRunAnyContext));
	operations.push_back(Operation("CreateXMLMetadata", "document", kSnpRunAnyContext));
	operations.push_back(Operation("ModifyXMLNode", "document", kSnpRunAnyContext));
	operations.push_back(Operation("ModifyXMLMetadata", "document", kSnpRunAnyContext));
	operations.push_back(Operation("DisplayXMLNode", "document", kSnpRunAnyContext));
	operations.push_back(Operation("DisplayXMLMetadata", "document", kSnpRunAnyContext));
	operations.push_back(Operation("RemoveXMLNodes", "document", kSnpRunAnyContext));
	operations.push_back(Operation("RemoveXMLMetadata", "document", kSnpRunAnyContext));
	return operations;
}

/*
*/
std::string _SnpRunnableXML::GetDefaultOperationName() const 
{
	return "CreateXMLNodes";
}

/*
*/
std::vector<std::string> _SnpRunnableXML::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("XML Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableXML::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	return selectionHelper.IsDocumentSelected();
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableXML::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpXML instance;
	if ("CreateXMLNodes" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateXMLNodes();
	}
	else if ("CreateXMLMetadata" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateXMLMetadata();
	}
	else if ("ModifyXMLNode" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyXMLNode();
	}
	else if ("ModifyXMLMetadata" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyXMLMetadata();
	}
	else if ("DisplayXMLNode" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayXMLNode(element);
	}
	else if ("DisplayXMLMetadata" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayXMLMetadata();
	}
	else if ("RemoveXMLNodes" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveXMLNodes();
	}
	else if ("RemoveXMLMetadata" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveXMLMetadata();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableXML instance_SnpRunnableXML;

// End SnpXML.cpp

