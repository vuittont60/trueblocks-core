#pragma once
/*-------------------------------------------------------------------------------------------
 * QuickBlocks - Decentralized, useful, and detailed data from Ethereum blockchains
 * Copyright (c) 2018 Great Hill Corporation (http://quickblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
#include <algorithm>
#include "basetypes.h"

namespace qblocks {

    //----------------------------------------------------------------------
    #define ARRAY_CHUNK_SIZE 100

    //-------------------------------------------------------------------------
    typedef int  (*SEARCHFUNC)    (const void *ob1, const void *ob2);
    typedef int  (*SORTINGFUNC)   (const void *ob1, const void *ob2);
    typedef int  (*DUPLICATEFUNC) (const void *ob1, const void *ob2);
    typedef bool (*APPLYFUNC)     (string_q& line, void *data);

    //----------------------------------------------------------------------
    template<class TYPE>
    class SFArrayBase {
    protected:
        size_t m_nCapacity;
        size_t m_nItems;
        TYPE  *m_Items;

    public:
        SFArrayBase(void);
        SFArrayBase(const SFArrayBase& cop);
        ~SFArrayBase(void);

        SFArrayBase& operator=(const SFArrayBase& cop);

              TYPE&  at        (size_t index);
        const TYPE&  operator[](size_t index) const;
              size_t capacity  (void) const { return m_nCapacity; }
              size_t size      (void) const { return m_nItems; }
              void   push_back (TYPE x);
              void   clear     (void);
              void   reserve   (size_t newSize);
              void   resize    (size_t newSize) { reserve(newSize); }

        void Sort(SORTINGFUNC func) { qsort(&m_Items[0], m_nItems, sizeof(TYPE), func); }
        TYPE *Find(const TYPE *key, SEARCHFUNC func) {
            // note: use the same function you would use to sort. Return <0, 0, or >0 if less, equal, greater
            return reinterpret_cast<TYPE*>(bsearch(key, &m_Items[0], m_nItems, sizeof(TYPE), func));
        }

    private:
        void checkSize(size_t sizeNeeded);
        void duplicate(const SFArrayBase& cop);
        void initialize(size_t cap, size_t count, TYPE *values);
    };

    //----------------------------------------------------------------------
    template<class TYPE>
    inline SFArrayBase<TYPE>::SFArrayBase(void) {
        initialize(0, 0, NULL);
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline SFArrayBase<TYPE>::SFArrayBase(const SFArrayBase<TYPE>& cop) {
        duplicate(cop);
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline SFArrayBase<TYPE>::~SFArrayBase(void) {
        clear();
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline SFArrayBase<TYPE>& SFArrayBase<TYPE>::operator=(const SFArrayBase<TYPE>& cop) {
        clear();
        duplicate(cop);
        return *this;
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::initialize(size_t cap, size_t count, TYPE *values) {
        m_nCapacity = cap;
        m_nItems = count;
        m_Items  = values;
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::clear(void) {
        if (m_Items)
            delete [] m_Items;
        initialize(0, 0, NULL);
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::duplicate(const SFArrayBase<TYPE>& cop) {
        checkSize(cop.capacity());
        for (size_t i = 0 ; i < cop.size() ; i++)
            m_Items[i] = cop.m_Items[i];
        m_nItems = cop.size();
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::reserve(size_t newSize) {
        checkSize(newSize);
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::checkSize(size_t sizeNeeded) {
        if (sizeNeeded < m_nCapacity)
            return;

        // The user is requesting access to an index that is past range. We need to resize the array.
        size_t newSize = max(m_nCapacity + ARRAY_CHUNK_SIZE, sizeNeeded);
        TYPE *newArray = new TYPE[newSize];
        if (m_nItems) {
            // If there are any values in the source copy them over
            for (size_t i = 0 ; i < m_nItems ; i++)
                newArray[i] = m_Items[i];
            // Then clear out the old array
            if (m_Items)
                delete [] m_Items;
            m_Items = NULL;
        }
        initialize(newSize, m_nItems, newArray);
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline TYPE& SFArrayBase<TYPE>::at(size_t index) {
        // TODO: This should definitly not grow the array. If we use this to grow the array, 
        // when we switch to a native vector, this will break
        checkSize(index);
        if (index >= m_nItems)
            m_nItems = index+1;
        ASSERT(m_Items && index >= 0 && index <= m_nCapacity && index <= m_nItems);
        return m_Items[index];
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline void SFArrayBase<TYPE>::push_back(TYPE x) {
        size_t index = size();
        checkSize(index);
        if (index >= m_nItems)
            m_nItems = index + 1;
        ASSERT(m_Items && index >= 0 && index <= m_nCapacity && index <= m_nItems);
        m_Items[index] = x;
    }

    //----------------------------------------------------------------------
    template<class TYPE>
    inline const TYPE& SFArrayBase<TYPE>::operator[](size_t index) const {
        // This is the const version which means it's a get which means we should not be expecting
        // the array to grow. Does not appear to protect against accessing outside range though.
        ASSERT(index >= 0 && index <= m_nItems);
        return m_Items[index];
    }

    struct xLISTPOS__ { int unused; };
    typedef xLISTPOS__* LISTPOS;

    //----------------------------------------------------------------------
    template<class TYPE>
    class SFListNode {
    private:
        // Don't allow default construction, copy construction or assignment
        SFListNode(void) { }
        SFListNode(const SFListNode& node) { }
        SFListNode& operator=(const SFListNode& node) { return *this; }

    public:
        SFListNode *m_Next;
        SFListNode *m_Prev;
        TYPE m_Data;

        explicit SFListNode(TYPE d) { m_Next = NULL; m_Prev = NULL;  m_Data = d; }
        ~SFListNode(void) { m_Next = NULL; m_Prev = NULL;  }

        SFListNode *Next(void) const { return m_Next; }
        SFListNode *Prev(void) const { return m_Prev; }
    };

    //----------------------------------------------------------------------
    template<class TYPE>
    class SFList {
    protected:
        size_t m_Count;
        SFListNode<TYPE> *m_Head;
        SFListNode<TYPE> *m_Tail;

    public:
        SFList(void);
        SFList(const SFList& l);
        ~SFList(void);

        SFList& operator=(const SFList& l);

        size_t size(void) const { return m_Count; }
        TYPE GetHead(void) const { return (TYPE)(m_Head->m_Data); }
        TYPE GetTail(void) const { return (TYPE)(m_Tail->m_Data); }

        LISTPOS GetHeadPosition (void) const { return (LISTPOS)m_Head; }
        LISTPOS GetTailPosition (void) const { return (LISTPOS)m_Tail; }

        void setHead(SFListNode<TYPE> *newHead) { m_Head = newHead; }
        void setTail(SFListNode<TYPE> *newTail) { m_Tail = newTail; }

        void AddToList(TYPE item) { AddTail(item); }
        bool empty(void) const { return (m_Head == NULL); }

        TYPE GetNext(LISTPOS& rPosition) const;
        TYPE GetPrev(LISTPOS& rPosition) const;
        LISTPOS Find(TYPE item) const;
        TYPE FindAt(TYPE item) const;
        TYPE FindAt(LISTPOS pos) const;

        void AddHead(TYPE item);
        void AddTail(TYPE item);
        void AddToList(const SFList& l);
        bool AddSorted(TYPE item, SORTINGFUNC sortFunc, DUPLICATEFUNC dupFunc = NULL);

        void InsertBefore(LISTPOS pos, TYPE item);
        void InsertAfter(LISTPOS pos, TYPE item);

        TYPE RemoveAt(LISTPOS pos);
        TYPE RemoveHead(void);
        TYPE RemoveTail(void);
        void RemoveAll(void);
    };

    //---------------------------------------------------------------------
    template<class TYPE>
    inline SFList<TYPE>::SFList(void) {
        m_Head  = NULL;
        m_Tail  = NULL;
        m_Count = 0;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    SFList<TYPE>::SFList(const SFList<TYPE>& l) {
        m_Head  = NULL;
        m_Tail  = NULL;
        m_Count = 0;

        LISTPOS pos = l.GetHeadPosition();
        while (pos) {
            TYPE ob = l.GetNext(pos);
            AddTail(ob);
        }
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline SFList<TYPE>::~SFList(void) {
        RemoveAll();
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    SFList<TYPE>& SFList<TYPE>::operator=(const SFList<TYPE>& l) {
        RemoveAll();

        LISTPOS pos = l.GetHeadPosition();
        while (pos) {
            TYPE ob = l.GetNext(pos);
            AddTail(ob);
        }
        return *this;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    TYPE SFList<TYPE>::FindAt(TYPE probe) const {
        LISTPOS pos = GetHeadPosition();
        while (pos) {
            TYPE ob = GetNext(pos);
            if (ob == probe)
                return ob;
        }
        return NULL;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    TYPE SFList<TYPE>::FindAt(LISTPOS probe) const {
        LISTPOS pos = GetHeadPosition();
        while (pos) {
            LISTPOS prev = pos;
            TYPE ob = GetNext(pos);
            if (prev == probe)
                return ob;
        }
        return NULL;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    LISTPOS SFList<TYPE>::Find(TYPE probe) const {
        LISTPOS pos = GetHeadPosition();
        while (pos) {
            LISTPOS last = pos;
            TYPE ob = GetNext(pos);
            if (ob == probe)
                return last;
        }
        return NULL;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline void SFList<TYPE>::AddHead(TYPE data) {
        SFListNode<TYPE> *node = new SFListNode<TYPE>(data);

        ASSERT(node);
        ASSERT(!m_Head || m_Head->m_Prev == NULL);
        ASSERT(!m_Tail || m_Tail->m_Next == NULL);

        node->m_Next = m_Head;
        node->m_Prev = NULL;

        if (!m_Head) {
            ASSERT(!m_Tail);
            m_Head = m_Tail = node;
        } else {
            ASSERT(m_Tail);
            m_Head->m_Prev = node;
            m_Head = node;
        }
        m_Count++;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline void SFList<TYPE>::AddTail(TYPE data) {
        SFListNode<TYPE> *node = new SFListNode<TYPE>(data);

        ASSERT(node);
        ASSERT(!m_Head || m_Head->m_Prev == NULL);
        ASSERT(!m_Tail || m_Tail->m_Next == NULL);

        node->m_Next = NULL;
        node->m_Prev = m_Tail;

        if (!m_Head) {
            ASSERT(!m_Tail);
            m_Head = m_Tail = node;
        } else {
            ASSERT(m_Tail);
            m_Tail->m_Next = node;
            m_Tail = node;
        }
        m_Count++;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline void SFList<TYPE>::InsertBefore(LISTPOS pos, TYPE data) {
        SFListNode<TYPE> *node = new SFListNode<TYPE>(data);
        SFListNode<TYPE> *before = (SFListNode<TYPE> *)pos;

        ASSERT(node && before);

        node->m_Prev = before->m_Prev;
        node->m_Next = before;

        if (before->m_Prev)
            before->m_Prev->m_Next = node;
        before->m_Prev = node;

        ASSERT(m_Head && m_Tail);  // We would have used AddTail otherwise
        if (before == m_Head)
            m_Head = node;

        m_Count++;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline void SFList<TYPE>::InsertAfter(LISTPOS pos, TYPE data) {
        SFListNode<TYPE> *node  = new SFListNode<TYPE>(data);
        SFListNode<TYPE> *after = (SFListNode<TYPE> *)pos;

        ASSERT(node && after);

        node->m_Prev = after;
        node->m_Next = after->m_Next;

        if (after->m_Next)
            after->m_Next->m_Prev = node;
        after->m_Next = node;

        ASSERT(m_Head && m_Tail);  // We would have used AddTail otherwise
        if (after == m_Tail)
            m_Tail = node;

        m_Count++;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    void SFList<TYPE>::AddToList(const SFList<TYPE>& l) {
        LISTPOS pos = l.GetHeadPosition();
        while (pos) {
            TYPE ob = l.GetNext(pos);
            AddToList(ob);
        }
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline TYPE SFList<TYPE>::GetNext(LISTPOS& pos) const {
        SFListNode<TYPE> *node = (SFListNode<TYPE> *)pos;
        ASSERT(node);
        pos = (LISTPOS)((node->m_Next != m_Head) ? node->m_Next : NULL);

        return (TYPE)(node->m_Data);
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline TYPE SFList<TYPE>::GetPrev(LISTPOS& pos) const {
        SFListNode<TYPE> *node = (SFListNode<TYPE> *)pos;
        ASSERT(node);
        pos = (LISTPOS)((node->m_Prev != m_Tail) ? node->m_Prev : NULL);

        return (TYPE)(node->m_Data);
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline void SFList<TYPE>::RemoveAll(void) {
        SFListNode<TYPE> *node = m_Head;
        while (node) {
            SFListNode<TYPE> *n = ((node->m_Next != m_Head) ? node->m_Next : NULL);
            if (node == m_Head)
                m_Head = NULL;
            delete node;
            node = n;
        }

        m_Head      = NULL;
        m_Tail      = NULL;
        m_Count     = 0;
    }

    //---------------------------------------------------------------------
    template<class TYPE>
    inline TYPE SFList<TYPE>::RemoveAt(LISTPOS pos) {
        SFListNode<TYPE> *node = (SFListNode<TYPE> *)pos;

        ASSERT(node);
        ASSERT(!m_Head || m_Head->m_Prev == NULL);
        ASSERT(!m_Tail || m_Tail->m_Next == NULL);

        TYPE data = (TYPE)(node->m_Data);

        if (!m_Head) {
            ASSERT(!m_Tail);
            delete node;
            return data;
        }
        ASSERT(m_Tail);

        if (m_Head == node)
            m_Head = m_Head->m_Next;

        if (m_Tail == node)
            m_Tail = m_Tail->m_Prev;

        if (node->m_Prev)
            node->m_Prev->m_Next = node->m_Next;

        if (node->m_Next)
            node->m_Next->m_Prev = node->m_Prev;

        m_Count--;

        delete node;
        return data;
    }

    //---------------------------------------------------------------------
    // stack use
    template<class TYPE>
    inline TYPE SFList<TYPE>::RemoveHead(void) {
        return RemoveAt((LISTPOS)m_Head);
    }

    //---------------------------------------------------------------------
    // queue use
    template<class TYPE>
    inline TYPE SFList<TYPE>::RemoveTail(void) {
        return RemoveAt((LISTPOS)m_Tail);
    }

    //-----------------------------------------------------------------------------
    // return true of added, false otherwise so caller can free allocated memory if any
    template<class TYPE>
    inline bool SFList<TYPE>::AddSorted(TYPE item, SORTINGFUNC sortFunc, DUPLICATEFUNC dupFunc) {
        if (!item)
            return false;

        if (sortFunc) {
            // Sort it in (if told to)...
            LISTPOS ePos = GetHeadPosition();
            while (ePos) {
                LISTPOS lastPos = ePos;
                TYPE test = GetNext(ePos);

                bool isDup = dupFunc && (dupFunc)(item, test);
                if (isDup) {
                    // caller must free this memory or drop it
                    return false;
                }

                if ((sortFunc)(item, test) < 0) {
                    InsertBefore(lastPos, item);
                    return true;
                }
            }
        }

        // ...else just add it to the end
        AddToList(item);
        return true;
    }

    //-----------------------------------------------------------------------------------------
    inline int sortByStringValue(const void *rr1, const void *rr2) {
        string_q n1 = * reinterpret_cast<const string_q*>(rr1);
        string_q n2 = * reinterpret_cast<const string_q*>(rr2);
        return strcasecmp(n1.c_str(), n2.c_str());
    }
}  // namespace qblocks

